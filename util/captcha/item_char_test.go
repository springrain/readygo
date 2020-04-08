package captcha

import (
	"bytes"
	"fmt"
	"image/color"
	"reflect"
	"testing"

	"github.com/golang/freetype/truetype"
)

func TestNewItemChar(t *testing.T) {

	bgColor := color.RGBA{0, 0, 0, 0}
	got := newItemChar(48, 30, bgColor)
	got.drawHollowLine()
	got.drawText("ab1c")

	fmt.Println(got.encodeB64string())

}

func TestItemChar_drawHollowLine(t *testing.T) {
	tests := []struct {
		name string
		item *itemChar
		want *itemChar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.drawHollowLine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemChar.drawHollowLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemChar_drawSineLine(t *testing.T) {
	tests := []struct {
		name string
		item *itemChar
		want *itemChar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.drawSineLine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemChar.drawSineLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemChar_drawSlimLine(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		item *itemChar
		args args
		want *itemChar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.drawSlimLine(tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemChar.drawSlimLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemChar_drawBeeline(t *testing.T) {
	type args struct {
		point1    point
		point2    point
		lineColor color.RGBA
	}
	tests := []struct {
		name string
		item *itemChar
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.item.drawBeeline(tt.args.point1, tt.args.point2, tt.args.lineColor)
		})
	}
}

func TestItemChar_drawNoise(t *testing.T) {
	type args struct {
		noiseText string
		fonts     []*truetype.Font
	}
	tests := []struct {
		name    string
		item    *itemChar
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.drawNoise(tt.args.noiseText); (err != nil) != tt.wantErr {
				t.Errorf("ItemChar.drawNoise() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestItemChar_drawText(t *testing.T) {
	type args struct {
		text  string
		fonts []*truetype.Font
	}
	tests := []struct {
		name    string
		item    *itemChar
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.drawText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("ItemChar.drawText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestItemChar_BinaryEncoding(t *testing.T) {
	tests := []struct {
		name string
		item *itemChar
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.binaryEncoding(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemChar.BinaryEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemChar_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		item    *itemChar
		want    int64
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.item.writeTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemChar.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ItemChar.WriteTo() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ItemChar.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestItemChar_EncodeB64string(t *testing.T) {
	tests := []struct {
		name string
		item *itemChar
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.encodeB64string(); got != tt.want {
				t.Errorf("ItemChar.EncodeB64string() = %v, want %v", got, tt.want)
			}
		})
	}
}
