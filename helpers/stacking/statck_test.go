package stacking

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/simulot/immich-go/immich/metadata"
)

type asset struct {
	ID        string
	FileName  string
	DateTaken time.Time
}

func Test_Stack(t *testing.T) {
	tc := []struct {
		name  string
		input []asset
		want  []Stack
	}{
		{
			name: "no stack JPG+DNG",
			input: []asset{
				{ID: "1", FileName: "IMG_1234.JPG", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
				{ID: "2", FileName: "IMG_1234.DNG", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.45.00")},
			},
		},
		{
			name: "issue #67",
			input: []asset{
				{ID: "1", FileName: "IMG_5580.HEIC", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
				{ID: "2", FileName: "IMG_5580.MP4", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
			},
		},
		{
			name: "stack JPG+DNG",
			input: []asset{
				{ID: "1", FileName: "IMG_1234.JPG", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
				{ID: "2", FileName: "IMG_1234.DNG", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
			},

			want: []Stack{
				{
					CoverID: "1",
					IDs:     []string{"2"},
					Date:    metadata.TakeTimeFromName("2023-10-01 10.15.00"),
					Names:   []string{"IMG_1234.JPG", "IMG_1234.DNG"},
				},
			},
		},
		{
			name: "stack BURST",
			input: []asset{
				{ID: "1", FileName: "IMG_20231014_183244.jpg", DateTaken: metadata.TakeTimeFromName("IMG_20231014_183244.jpg")},
				{ID: "2", FileName: "IMG_20231014_183246_BURST001_COVER.jpg", DateTaken: metadata.TakeTimeFromName("IMG_20231014_183246_BURST001_COVER.jpg")},
				{ID: "3", FileName: "IMG_20231014_183246_BURST002.jpg", DateTaken: metadata.TakeTimeFromName("IMG_20231014_183246_BURST002.jpg")},
				{ID: "4", FileName: "IMG_20231014_183246_BURST003.jpg", DateTaken: metadata.TakeTimeFromName("IMG_20231014_183246_BURST003.jpg")},
			},
			want: []Stack{
				{
					CoverID: "2",
					IDs:     []string{"3", "4"},
					Date:    metadata.TakeTimeFromName("IMG_20231014_183246_BURST001_COVER.jpg"),
					Names:   []string{"IMG_20231014_183246_BURST001_COVER.jpg", "IMG_20231014_183246_BURST002.jpg", "IMG_20231014_183246_BURST003.jpg"},
				},
			},
		},

		{
			name: "stack JPG+CR3",
			input: []asset{
				{ID: "1", FileName: "3H2A0018.CR3", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
				{ID: "2", FileName: "3H2A0018.JPG", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
				{ID: "3", FileName: "3H2A0019.CR3", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
				{ID: "4", FileName: "3H2A0019.JPG", DateTaken: metadata.TakeTimeFromName("2023-10-01 10.15.00")},
			},
			want: []Stack{
				{
					CoverID: "2",
					IDs:     []string{"1"},
					Date:    metadata.TakeTimeFromName("2023-10-01 10.15.00"),
					Names:   []string{"3H2A0018.CR3", "3H2A0018.JPG"},
				},
				{
					CoverID: "4",
					IDs:     []string{"3"},
					Date:    metadata.TakeTimeFromName("2023-10-01 10.15.00"),
					Names:   []string{"3H2A0019.CR3", "3H2A0019.JPG"},
				},
			},
		},
		{
			name: "issue #12 example1",
			input: []asset{
				{ID: "1", FileName: "PXL_20231026_210642603.dng", DateTaken: metadata.TakeTimeFromName("PXL_20231026_210642603.dng")},
				{ID: "2", FileName: "PXL_20231026_210642603.jpg", DateTaken: metadata.TakeTimeFromName("PXL_20231026_210642603.jpg")},
			},
			want: []Stack{
				{
					CoverID: "2",
					IDs:     []string{"1"},
					Date:    metadata.TakeTimeFromName("PXL_20231026_210642603.dng"),
					Names:   []string{"PXL_20231026_210642603.dng", "PXL_20231026_210642603.jpg"},
				},
			},
		},
		{
			name: "issue #12 example 2",
			input: []asset{
				{ID: "3", FileName: "20231026_205755225.dng", DateTaken: metadata.TakeTimeFromName("20231026_205755225.dng")},
				{ID: "4", FileName: "20231026_205755225.MP.jpg", DateTaken: metadata.TakeTimeFromName("20231026_205755225.MP.jpg")},
			},
			want: []Stack{
				{
					CoverID: "4",
					IDs:     []string{"3"},
					Date:    metadata.TakeTimeFromName("20231026_205755225.MP.jpg"),
					Names:   []string{"20231026_205755225.dng", "20231026_205755225.MP.jpg"},
				},
			},
		},
		{
			name: "issue #12 example 3",
			input: []asset{
				{ID: "3", FileName: "20231026_205755225.dng", DateTaken: metadata.TakeTimeFromName("20231026_205755225.dng")},
				{ID: "4", FileName: "20231026_205755225.MP.jpg", DateTaken: metadata.TakeTimeFromName("20231026_205755225.MP.jpg")},
				{ID: "5", FileName: "PXL_20231207_032111247.RAW-02.ORIGINAL.dng", DateTaken: metadata.TakeTimeFromName("PXL_20231207_032111247.RAW-02.ORIGINAL.dng")},
				{ID: "6", FileName: "PXL_20231207_032111247.RAW-01.COVER.jpg", DateTaken: metadata.TakeTimeFromName("PXL_20231207_032111247.RAW-01.COVER.jpg")},
				{ID: "7", FileName: "PXL_20231207_032108788.RAW-02.ORIGINAL.dng", DateTaken: metadata.TakeTimeFromName("PXL_20231207_032108788.RAW-02.ORIGINAL.dng")},
				{ID: "8", FileName: "PXL_20231207_032108788.RAW-01.MP.COVER.jpg", DateTaken: metadata.TakeTimeFromName("PXL_20231207_032108788.RAW-01.MP.COVER.jpg")},
			},
			want: []Stack{
				{
					CoverID: "4",
					IDs:     []string{"3"},
					Date:    metadata.TakeTimeFromName("20231026_205755225.dng"),
					Names:   []string{"20231026_205755225.dng", "20231026_205755225.MP.jpg"},
				},
				{
					CoverID: "6",
					IDs:     []string{"5"},
					Date:    metadata.TakeTimeFromName("PXL_20231207_032111247.RAW-02.ORIGINAL.dng"),
					Names:   []string{"PXL_20231207_032111247.RAW-02.ORIGINAL.dng", "PXL_20231207_032111247.RAW-01.COVER.jpg"},
				},
				{
					CoverID: "8",
					IDs:     []string{"7"},
					Date:    metadata.TakeTimeFromName("PXL_20231207_032108788.RAW-02.ORIGINAL.dng"),
					Names:   []string{"PXL_20231207_032108788.RAW-02.ORIGINAL.dng", "PXL_20231207_032108788.RAW-01.MP.COVER.jpg"},
				},
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewStackBuilder()
			for _, a := range tt.input {
				sb.ProcessAsset(a.ID, a.FileName, a.DateTaken)
			}

			got := sb.Stacks()
			sort.Slice(got, func(i, j int) bool {
				return got[i].CoverID < got[j].CoverID
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("difference\n")
				pretty.Ldiff(t, tt.want, got)
			}
		})

	}
}
