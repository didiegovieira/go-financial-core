package runtime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleSize() {
	sizes := []Size{
		3 * Byte,
		1000 * Byte,
		1024 * Byte,
		1536 * Byte,
		2048 * Byte,

		3 * Megabyte,
		1000 * Megabyte,
		1024 * Megabyte,
		1536 * Megabyte,
		2048 * Megabyte,

		3 * Terabyte,
		1000 * Terabyte,
		1024 * Terabyte,
		1536 * Terabyte,
		2048 * Terabyte,
	}

	for _, s := range sizes {
		fmt.Printf("%.1f\n", s)
	}

	// Output:
	// 3.0B
	// 1000.0B
	// 1.0KB
	// 1.5KB
	// 2.0KB
	// 3.0MB
	// 1000.0MB
	// 1.0GB
	// 1.5GB
	// 2.0GB
	// 3.0TB
	// 1000.0TB
	// 1024.0TB
	// 1536.0TB
	// 2048.0TB
}

func TestExampleSizeFormat(t *testing.T) {
	size := 100 * Megabyte

	fmt.Printf("   d = %d\n", size)
	fmt.Printf("   s = %s\n", size)
	fmt.Printf("   q = %q\n", size)
	fmt.Printf("   v = %v\n", size)
	fmt.Printf("   f = %f\n", size)
	fmt.Printf(" .2f = %.2f\n", size)

	fmt.Printf("  8d = %8d\n", size)
	fmt.Printf("  8s = %8s\n", size)
	fmt.Printf("  8q = %8q\n", size)
	fmt.Printf("  8v = %8v\n", size)
	fmt.Printf("  8f = %8f\n", size)
	fmt.Printf("8.2f = %8.2f\n", size)

	// Output:
	//    d = 100MB
	//    s = 100MB
	//    q = "100MB"
	//    v = 100MB
	//    f = 100MB
	//  .2f = 100.00MB

	//   8d =    100MB
	//   8s =    100MB
	//   8q = 	"100MB"
	//   8v =    100MB
	//   8f =    100MB
	// 8.2f = 100.00MB

	assert.Equal(t, 1, 1)
}

func TestSizeFormatCanFormatWithEveryFlag(t *testing.T) {
	size := 100 * Megabyte

	cases := []struct {
		format string
		want   string
	}{
		{format: "%d", want: `100MB`},
		{format: "%s", want: `100MB`},
		{format: "%q", want: `"100MB"`},
		{format: "%v", want: `100MB`},
		{format: "%f", want: `100MB`},
		{format: "%.2f", want: `100.00MB`},
		{format: "%8d", want: `   100MB`},
		{format: "%8s", want: `   100MB`},
		{format: "%8q", want: `"   100MB"`},
		{format: "%8v", want: `   100MB`},
		{format: "%8f", want: `   100MB`},
		{format: "%8.2f", want: `100.00MB`},
	}

	for _, tc := range cases {
		t.Run(tc.format, func(t *testing.T) {
			got := fmt.Sprintf(tc.format, size)
			assert.Equal(t, got, tc.want)
		})
	}
}

func TestSizeFormatWorksWithFractionals(t *testing.T) {
	format := "%.2f"

	cases := []struct {
		want string
		size Size
	}{
		{want: "1.00KB", size: 1024 * Byte},
		{want: "1.25KB", size: 1280 * Byte},
		{want: "1.50KB", size: 1536 * Byte},
		{want: "1.75KB", size: 1792 * Byte},

		{want: "2.00KB", size: 2048 * Byte},
		{want: "2.25KB", size: 2304 * Byte},
		{want: "2.50KB", size: 2560 * Byte},
		{want: "2.75KB", size: 2816 * Byte},

		{want: "3.00KB", size: 3072 * Byte},
		{want: "3.25KB", size: 3328 * Byte},
		{want: "3.50KB", size: 3584 * Byte},
		{want: "3.75KB", size: 3840 * Byte},
	}

	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			got := fmt.Sprintf(format, tc.size)
			assert.Equal(t, got, tc.want)
		})
	}
}

func TestParseSizeCanParseCorrectSizes(t *testing.T) {
	cases := []struct {
		arg  string
		want Size
	}{
		{arg: "0", want: 0},
		{arg: "0B", want: 0},
		{arg: "0KB", want: 0},
		{arg: "0MB", want: 0},
		{arg: "0GB", want: 0},
		{arg: "0TB", want: 0},

		{arg: "256B", want: 256 * Byte},
		{arg: "1536B", want: 1536 * Byte},
		{arg: "2048B", want: 2048 * Byte},

		{arg: "0.25KB", want: 256 * Byte},
		{arg: "1.5KB", want: 1536 * Byte},
		{arg: "2KB", want: 2048 * Byte},

		{arg: "256KB", want: 256 * Kilobyte},
		{arg: "1536KB", want: 1536 * Kilobyte},
		{arg: "2048KB", want: 2048 * Kilobyte},

		{arg: "0.25MB", want: 256 * Kilobyte},
		{arg: "1.5MB", want: 1536 * Kilobyte},
		{arg: "2MB", want: 2048 * Kilobyte},

		{arg: "256MB", want: 256 * Megabyte},
		{arg: "1536MB", want: 1536 * Megabyte},
		{arg: "2048MB", want: 2048 * Megabyte},

		{arg: "0.25GB", want: 256 * Megabyte},
		{arg: "1.5GB", want: 1536 * Megabyte},
		{arg: "2GB", want: 2048 * Megabyte},

		{arg: "256GB", want: 256 * Gigabyte},
		{arg: "1536GB", want: 1536 * Gigabyte},
		{arg: "2048GB", want: 2048 * Gigabyte},

		{arg: "0.25TB", want: 256 * Gigabyte},
		{arg: "1.5TB", want: 1536 * Gigabyte},
		{arg: "2TB", want: 2048 * Gigabyte},

		{arg: "256TB", want: 256 * Terabyte},
		{arg: "1536TB", want: 1536 * Terabyte},
		{arg: "2048TB", want: 2048 * Terabyte},
	}

	for _, tc := range cases {
		t.Run(tc.arg, func(t *testing.T) {
			got, err := ParseSize(tc.arg)
			assert.Nil(t, err)
			assert.Equal(t, got, tc.want)
		})
	}
}

func TestParseSizeCanBeRevertTheOutputOfString(t *testing.T) {
	cases := []struct {
		want Size
	}{
		{want: 42 * Byte},
		{want: 1536 * Byte},
		{want: 42 * Kilobyte},
		{want: 1536 * Kilobyte},
		{want: 42 * Megabyte},
		{want: 1536 * Megabyte},
		{want: 42 * Gigabyte},
		{want: 1536 * Gigabyte},
		{want: 42 * Terabyte},
		{want: 1536 * Terabyte},
	}

	for _, tc := range cases {
		t.Run(tc.want.String(), func(t *testing.T) {
			got, err := ParseSize(tc.want.String())
			assert.Nil(t, err)
			assert.Equal(t, got, tc.want)
		})
	}
}
