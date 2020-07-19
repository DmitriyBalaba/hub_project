package payload

import (
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
)

const MultipartMaxMemory = 102400

func MultipartBody(dto, model interface{}, key string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			getCopyByType := func(t interface{}) interface{} {
				return reflect.New(reflect.TypeOf(t)).Interface()
			}

			dto, model := getCopyByType(dto), getCopyByType(model)

			if err := r.ParseMultipartForm(MultipartMaxMemory); err != nil {
				ErrorResponse(w, WrapBadRequestErr(err))
				return
			}

			multipartVal, ok := r.MultipartForm.Value[key]
			if !ok {
				ErrorResponse(w, NewBadRequestErr("%s not found in multipart", key))
				return
			}

			// TODO: handle array values instead of taking the first value

			readCloser := ioutil.NopCloser(strings.NewReader(multipartVal[0]))

			if err := getValid(readCloser, dto); err != nil {
				ErrorResponse(w, WrapBadRequestErr(err))
				return
			}

			if err := copyStructFields(dto, model); err != nil {
				ErrorResponse(w, err)
				return
			}

			next.ServeHTTP(w, withDtoAndModel(r, dto, model))
		})
	}
}

func OpenAllMultipart(r *http.Request, key string) (files []io.ReadCloser, filenames []string, err error) {
	err = r.ParseMultipartForm(MultipartMaxMemory)
	if err != nil {
		err = WrapBadRequestErr(errors.WithStack(err))
		return
	}

	headers, ok := r.MultipartForm.File[key]
	if !ok {
		return
	}

	filenames = make([]string, 0, len(headers))
	files = make([]io.ReadCloser, 0, len(headers))
	defer func() {
		if err != nil {
			for i := range files {
				_ = files[i].Close()
			}
		}
	}()

	var file multipart.File

	for _, h := range headers {
		file, err = h.Open()
		if err != nil {
			err = errors.WithStack(err)
			return
		}

		files = append(files, file)
		filenames = append(filenames, h.Filename)
	}

	return files, filenames, nil
}
