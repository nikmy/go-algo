package buffer

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBuffer_ZeroValue(t *testing.T) {
	type testcase struct {
		name string
		meth func(buf *Buffer)
	}

	tests := []testcase{
		{
			name: "Resize",
			meth: func(buf *Buffer) { buf.Resize(256) },
		},
		{
			name: "Free",
			meth: (*Buffer).Free,
		},
		{
			name: "ReadAll",
			meth: func(buf *Buffer) {
				r := NewMockReader(gomock.NewController(t))
				r.EXPECT().Close().Return(nil)
				r.EXPECT().Read(gomock.Any()).Do(func(b []byte) {}).Return(0, io.EOF)
				_ = buf.ReadAll(r)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf Buffer
			require.NotPanics(t, func() { tt.meth(&buf) })
		})
	}
}

func TestBuffer_ReadAll(t *testing.T) {
	const mockData = "0123456789"
	t.Run("single pass", func(t *testing.T) {
		buf := new(Buffer)
		buf.data = make([]byte, 10)

		r := NewMockReader(gomock.NewController(t))
		r.EXPECT().Close().Return(nil)
		r.EXPECT().
			Read(gomock.Any()).
			Do(func(b []byte) { copy(b, mockData) }).
			Return(10, io.EOF)

		err := buf.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, mockData, string(buf.Data()))
	})

	t.Run("read with error", func(t *testing.T) {
		buf := new(Buffer)
		buf.data = make([]byte, 10)

		r := NewMockReader(gomock.NewController(t))
		r.EXPECT().Close().Return(nil)
		r.EXPECT().
			Read(gomock.Any()).
			Return(0, os.ErrClosed)

		err := buf.ReadAll(r)
		require.ErrorIs(t, err, os.ErrClosed)
	})

	t.Run("grow needed", func(t *testing.T) {
		buf := new(Buffer)
		buf.data = make([]byte, 4, 8)

		r := NewMockReader(gomock.NewController(t))
		r.EXPECT().Close().Return(nil)
		r.EXPECT().
			Read(gomock.Any()).
			Do(func(b []byte) { copy(b, mockData) }).
			Return(8, nil)

		r.EXPECT().
			Read(gomock.Any()).
			Do(func(b []byte) { copy(b, mockData[8:]) }).
			Return(2, io.EOF)

		err := buf.ReadAll(r)
		require.ErrorIs(t, err, nil)
	})
}

func TestBuffer_Resize(t *testing.T) {
	type state struct {
		len int
		cap int
	}

	type testcase struct {
		name       string
		current    state
		desired    int
		want       state
		capAtLeast bool
	}

	tests := []testcase{
		{
			name: "grow within cap",
			current: state{
				len: 42,
				cap: 64,
			},
			desired: 50,
			want: state{
				len: 50,
				cap: 64,
			},
		},
		{
			name: "shrink",
			current: state{
				len: 50,
				cap: 64,
			},
			desired: 30,
			want: state{
				len: 30,
				cap: 64,
			},
		},
		{
			name: "keep size",
			current: state{
				len: 42,
				cap: 64,
			},
			desired: 42,
			want: state{
				len: 42,
				cap: 64,
			},
		},
		{
			name: "grow full",
			current: state{
				len: 256,
				cap: 256,
			},
			desired: 333,
			want: state{
				len: 333,
				cap: 333,
			},
			capAtLeast: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Buffer{
				data: make([]byte, tt.current.len, tt.current.cap),
			}

			require.NotPanics(t, func() {
				b.Resize(tt.desired)
			})

			require.Equal(t, tt.want.len, len(b.data))
			if tt.capAtLeast {
				require.LessOrEqual(t, tt.want.cap, cap(b.data))
			} else {
				require.Equal(t, tt.want.cap, cap(b.data))
			}
		})
	}
}
