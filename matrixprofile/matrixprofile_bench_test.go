package matrixprofile

import (
	"testing"
)

func setupData() []float64 {
	line := Line(0, 0, 512)
	ext := Line(0, 100, len(line)/2)
	ext2 := Line(0, 600, len(line)/2)
	sig := append(line, ext...)
	sig = append(sig, ext2...)
	noise := Noise(10, len(sig))
	sig = SigAdd(sig, noise)

	return sig
}

func BenchmarkZNormalize(b *testing.B) {
	b.ReportAllocs()
	sig := setupData()
	q := sig[:32]
	var err error
	var qnorm []float64
	for i := 0; i < b.N; i++ {
		qnorm, err = zNormalize(q)
		if err != nil {
			b.Error(err)
		}
		if len(qnorm) < 1 {
			b.Error("expected at least one value from z-normalizing a timeseries")
		}
	}
}

func BenchmarkMovstd(b *testing.B) {
	b.ReportAllocs()
	sig := setupData()
	var err error
	var std []float64
	for i := 0; i < b.N; i++ {
		std, err = movstd(sig, 32)
		if err != nil {
			b.Error(err)
		}
		if len(std) < 1 {
			b.Error("expected at least one value from moving standard deviation of a timeseries")
		}
	}
}

func BenchmarkCrossCorrelate(b *testing.B) {
	b.ReportAllocs()
	sig := setupData()
	q := sig[:32]
	var err error
	var cc []float64

	mp, err := New(q, sig, 32)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		cc, err = mp.crossCorrelate(q)
		if err != nil {
			b.Error(err)
		}
		if len(cc) < 1 {
			b.Error("expected at least one value from cross correlation of a timeseries")
		}
	}
}

func BenchmarkMass(b *testing.B) {
	b.ReportAllocs()
	sig := setupData()
	var err error
	var q, mprof []float64

	mp, err := New(sig, sig, 32)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		q = sig[:32]
		mprof, err = mp.mass(q)
		if err != nil {
			b.Error(err)
		}
		if len(mprof) < 1 {
			b.Error("expected at least one value from matrix profile")
		}
	}
}

func BenchmarkDistanceProfile(b *testing.B) {
	b.ReportAllocs()
	sig := setupData()
	var err error
	var mprof []float64

	mp, err := New(sig, nil, 32)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		mprof, err = mp.distanceProfile(0)
		if err != nil {
			b.Error(err)
		}
		if len(mprof) < 1 {
			b.Error("expected at least one value from matrix profile")
		}
	}
}

func BenchmarkStmp(b *testing.B) {
	sig := setupData()

	benchmarks := []struct {
		name string
		m    int
	}{
		{"m16", 16},
		{"m32", 32},
		{"m64", 64},
		{"m128", 128},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()

			mp, err := New(sig, nil, 32)
			if err != nil {
				b.Error(err)
			}

			for i := 0; i < b.N; i++ {
				err = mp.Stmp()
				if err != nil {
					b.Error(err)
				}
				if len(mp.MP) < 1 || len(mp.Idx) < 1 {
					b.Error("expected at least one value from matrix profile and matrix profile index")
				}
			}
		})
	}
}
