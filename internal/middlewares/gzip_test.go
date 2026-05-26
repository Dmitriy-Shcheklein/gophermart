package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipWriter(t *testing.T) {
	t.Run(
		"Write делегирует базовому ResponseWriter", func(t *testing.T) {
			underlyingWriter := &mockResponseWriter{
				status: http.StatusOK,
				size:   0,
			}

			gzWriter := gzipWriter{
				ResponseWriter: underlyingWriter,
				Writer:         &bytes.Buffer{},
			}

			testData := []byte("test data")
			n, err := gzWriter.Write(testData)

			require.NoError(t, err)
			assert.Equal(t, len(testData), n)
		},
	)

	t.Run(
		"Write возвращает корректные значения", func(t *testing.T) {
			writerBuffer := bytes.NewBufferString("already written")
			underlyingWriter := &mockResponseWriter{
				status: http.StatusOK,
				size:   0,
			}

			gzWriter := gzipWriter{
				ResponseWriter: underlyingWriter,
				Writer:         writerBuffer,
			}

			testData := []byte("new data")
			n, err := gzWriter.Write(testData)

			require.NoError(t, err)
			assert.Equal(t, len(testData), n)
		},
	)

	t.Run(
		"Write с ошибкой базового Writer", func(t *testing.T) {
			errorWriter := &errorWriter{err: assert.AnError}
			underlyingWriter := &mockResponseWriter{
				status: http.StatusOK,
				size:   0,
			}

			gzWriter := gzipWriter{
				ResponseWriter: underlyingWriter,
				Writer:         errorWriter,
			}

			testData := []byte("test data")
			n, err := gzWriter.Write(testData)

			assert.Error(t, err)
			assert.Equal(t, 0, n)
		},
	)
}

func TestDecompressRequest(t *testing.T) {
	var (
		logger      *zerolog.Logger
		nextHandler http.Handler
		r           *http.Request
		w           *httptest.ResponseRecorder

		mw *Middleware
	)

	setup := func(t *testing.T) {
		nopLogger := zerolog.Nop()
		logger = &nopLogger
		w = httptest.NewRecorder()
		nextHandler = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("response"))
			},
		)

		mw, _ = New(logger)
	}

	t.Run(
		"Успешная декомпрессия", func(t *testing.T) {
			setup(t)

			originalData := "original data"
			compressedData, err := compressData(originalData)
			require.NoError(t, err)

			r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(compressedData))
			r.Header.Set("Content-Encoding", "gzip")
			r.ContentLength = int64(len(compressedData))

			handler := mw.WithGzip(nextHandler)
			handler.ServeHTTP(w, r)

			decompressed, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Equal(t, originalData, string(decompressed))
			assert.Equal(t, int64(len(originalData)), r.ContentLength)
			assert.Empty(t, r.Header.Get("Content-Encoding"))
			assert.Equal(t, http.StatusOK, w.Code)
		},
	)

	t.Run(
		"Ошибка создания gzip.NewReader", func(t *testing.T) {
			setup(t)
			r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid gzip data"))
			r.Header.Set("Content-Encoding", "gzip")

			handler := mw.WithGzip(nextHandler)
			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "error while read gzip")
		},
	)

	t.Run(
		"Ошибка чтения из gzip reader", func(t *testing.T) {
			setup(t)
			r = httptest.NewRequest(http.MethodPost, "/", &errorReader{err: assert.AnError})
			r.Header.Set("Content-Encoding", "gzip")

			handler := mw.WithGzip(nextHandler)
			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "error while read gzip")
		},
	)

	t.Run(
		"Проверка удаления Content-Encoding", func(t *testing.T) {
			setup(t)
			originalData := "test data"
			compressedData, err := compressData(originalData)
			require.NoError(t, err)

			r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(compressedData))
			r.Header.Set("Content-Encoding", "gzip")

			handler := mw.WithGzip(nextHandler)
			handler.ServeHTTP(w, r)

			contentEncoding := r.Header.Get("Content-Encoding")
			assert.Empty(t, contentEncoding, "Content-Encoding должен быть удален")
		},
	)

	t.Run(
		"Проверка обновления ContentLength", func(t *testing.T) {
			setup(t)
			originalData := "data for length check"
			compressedData, err := compressData(originalData)
			require.NoError(t, err)

			r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(compressedData))
			r.Header.Set("Content-Encoding", "gzip")

			handler := mw.WithGzip(nextHandler)
			handler.ServeHTTP(w, r)

			assert.Equal(t, int64(len(originalData)), r.ContentLength)
		},
	)
}

func TestWithGzipMiddleware(t *testing.T) {
	var (
		logger      *zerolog.Logger
		nextHandler http.Handler
		r           *http.Request
		w           *httptest.ResponseRecorder

		mw *Middleware
	)

	setup := func(t *testing.T) {
		nopLogger := zerolog.Nop()
		logger = &nopLogger
		nextHandler = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("response"))
			},
		)

		mw, _ = New(logger)
	}

	t.Run(
		"Нет gzip заголовков", func(t *testing.T) {
			setup(t)

			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/", nil)

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "response", w.Body.String())
			assert.Empty(t, w.Header().Get("Content-Encoding"))
		},
	)

	t.Run(
		"Только Accept-Encoding - невалидный Content-Type", func(t *testing.T) {
			setup(t)
			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Type", "application/xml")

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "response", w.Body.String())
			assert.Empty(t, w.Header().Get("Content-Encoding"))
		},
	)

	t.Run(
		"Только Accept-Encoding - валидный Content-Type (application/json)", func(t *testing.T) {
			setup(t)

			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
			assert.NotEmpty(t, w.Body.Bytes())
		},
	)

	t.Run(
		"Только Accept-Encoding - валидный Content-Type (text/plain)", func(t *testing.T) {
			setup(t)

			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Type", "text/plain")

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
			assert.NotEmpty(t, w.Body.Bytes())
		},
	)

	t.Run(
		"Только Accept-Encoding - валидный Content-Type - проверка сжатия", func(t *testing.T) {
			setup(t)

			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

			compressedData := w.Body.Bytes()

			reader, err := gzip.NewReader(bytes.NewReader(compressedData))
			require.NoError(t, err)
			defer func() { _ = reader.Close() }()

			decompressed, err := io.ReadAll(reader)
			require.NoError(t, err)
			assert.Equal(t, "response", string(decompressed))
		},
	)

	t.Run(
		"Только Content-Encoding - успешная декомпрессия", func(t *testing.T) {
			setup(t)
			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()

			originalData := "compressed request data"
			compressedData, err := compressData(originalData)
			require.NoError(t, err)

			r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(compressedData))
			r.Header.Set("Content-Encoding", "gzip")
			r.Header.Set("Content-Type", "text/plain")

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Empty(t, w.Header().Get("Content-Encoding"))
		},
	)

	t.Run(
		"Только Content-Encoding - ошибка декомпрессии", func(t *testing.T) {
			setup(t)

			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid gzip"))
			r.Header.Set("Content-Encoding", "gzip")

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), "error while read gzip")
		},
	)

	t.Run(
		"Одновременный Accept-Encoding и Content-Encoding", func(t *testing.T) {
			setup(t)

			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()

			originalData := "compressed request"
			compressedRequest, err := compressData(originalData)
			require.NoError(t, err)

			r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(compressedRequest))
			r.Header.Set("Content-Encoding", "gzip")
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "response")
			assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
		},
	)

	t.Run(
		"Проверка gzipWriter.Write интерфейса", func(t *testing.T) {
			setup(t)

			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Accept-Encoding", "gzip")
			r.Header.Set("Content-Type", "text/plain")

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))

			reader, err := gzip.NewReader(bytes.NewReader(w.Body.Bytes()))
			require.NoError(t, err)
			defer func() { _ = reader.Close() }()

			receivedData, err := io.ReadAll(reader)
			require.NoError(t, err)
			assert.Equal(t, "response", string(receivedData))
		},
	)

	t.Run(
		"Проверка Content-Encoding установлен только при Accept-Encoding", func(t *testing.T) {
			setup(t)

			handler := mw.WithGzip(nextHandler)
			w = httptest.NewRecorder()
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Content-Encoding", "gzip")
			r.Header.Set("Content-Type", "application/json")

			originalData := "compressed request"
			compressedData, err := compressData(originalData)
			require.NoError(t, err)

			r.Body = io.NopCloser(bytes.NewReader(compressedData))

			handler.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Empty(
				t, w.Header().Get("Content-Encoding"),
				"Content-Encoding должен быть установлен только для response",
			)
		},
	)
}

func compressData(data string) ([]byte, error) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := gzWriter.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type errorWriter struct {
	err error
}

func (w *errorWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

type errorReader struct {
	err error
}

func (r *errorReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

type mockResponseWriter struct {
	writeResult struct {
		n   int
		err error
	}
	status int
	size   int
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	m.size += len(b)
	return m.writeResult.n, m.writeResult.err
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.status = statusCode
}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}
