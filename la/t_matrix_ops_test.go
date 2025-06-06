// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"math"
	"testing"

	"github.com/biotinker/gosl/chk"
	"github.com/biotinker/gosl/io"
	"github.com/biotinker/gosl/utl"
)

func checkIdentity(tst *testing.T, k string, a *Matrix, tol float64) {
	for i := 0; i < a.M; i++ {
		for j := 0; j < a.N; j++ {
			if i == j {
				diff := math.Abs(a.Get(i, j) - 1.0)
				if diff > tol {
					io.Pf("matrix %q is not diagonal. diff=%g\n", k, diff)
					tst.Errorf("matrix %q is not diagonal. diff=%g\n", k, diff)
					return
				}
			} else {
				diff := math.Abs(a.Get(i, j))
				if diff > tol {
					io.Pf("matrix %q is not diagonal. diff=%g\n", k, diff)
					tst.Errorf("matrix %q is not diagonal. diff=%g\n", k, diff)
					return
				}
			}
		}
	}
}

func checkAiSmall(tst *testing.T, k string, a *Matrix, zeroDet float64,
	correctAi *Matrix, correctDet, tolDet, tolAi, tolI float64) {

	// compute inverse and determinant
	ai := NewMatrix(a.M, a.N)
	det := MatInvSmall(ai, a, zeroDet)

	// check inverse and determinant
	chk.Float64(tst, "det("+k+")", tolDet, det, correctDet)
	chk.Array(tst, k+"i", tolAi, ai.Data, correctAi.Data)

	// check a⋅ai = I
	aai := NewMatrix(a.M, a.N)
	MatMatMul(aai, 1, a, ai)
	checkIdentity(tst, k+"⋅"+k+"i", aai, tolI)
}

func checkAi(tst *testing.T, k string, a *Matrix,
	correctAi *Matrix, correctDet, tolDet, tolAi, tolI float64) {

	// compute inverse
	ai := NewMatrix(a.N, a.M)
	det := MatInv(ai, a, true)

	// check determinant
	if a.M == a.N {
		chk.AnaNum(tst, "det("+k+") ", tolDet, det, correctDet, chk.Verbose)
		ddet := a.Det()
		chk.AnaNum(tst, k+".Det()", tolDet, ddet, correctDet, chk.Verbose)
	}

	// compare inverse
	if correctAi != nil {
		chk.Array(tst, k+"i", tolAi, ai.Data, correctAi.Data)
	}

	// multiply ai by a
	aai := NewMatrix(a.M, a.M)
	MatMatMul(aai, 1, a, ai)

	// square: check a⋅ai = I
	if a.M == a.N {
		checkIdentity(tst, k+"⋅"+k+"i", aai, tolI)
		return
	}

	// rectangular: check a⋅ai⋅a = a
	aaia := NewMatrix(a.M, a.N)
	MatMatMul(aaia, 1, aai, a)
	chk.Array(tst, k+"⋅"+k+"i⋅"+k, tolI, aaia.Data, a.Data)
}

func checkSvd(tst *testing.T, k string, a *Matrix, correctS []float64, correctU, correctVt *Matrix, tolS, tolU, tolVt, tolUsv float64) {

	// compute SVD
	s := make([]float64, utl.Imin(a.M, a.N))
	u := NewMatrix(a.M, a.M)
	vt := NewMatrix(a.N, a.N)
	MatSvd(s, u, vt, a, true)

	// compare results
	if correctS != nil {
		chk.Array(tst, k+": s", tolS, s, correctS)
	}
	if correctU != nil {
		chk.Array(tst, k+": u", tolU, u.Data, correctU.Data)
	}
	if correctVt != nil {
		chk.Array(tst, k+": vt", tolVt, vt.Data, correctVt.Data)
	}

	// check u⋅s⋅vt
	usv := NewMatrix(a.M, a.N)
	for i := 0; i < a.M; i++ {
		for j := 0; j < a.N; j++ {
			for k := 0; k < len(s); k++ {
				usv.Add(i, j, u.Get(i, k)*s[k]*vt.Get(k, j))
			}
		}
	}
	chk.Array(tst, k+": u⋅s⋅vt", tolUsv, usv.Data, a.Data)
}

func TestMatInv01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MatInv01. 1x1, 2x2, 3x3 matrices")

	// 1 x 1 matrix
	io.Pf("\n----------------------------------(1 x 1)-----------------------------------\n")
	a := NewMatrixDeep2([][]float64{
		{2.0},
	})
	ai := NewMatrixDeep2([][]float64{
		{0.5},
	})
	checkAiSmall(tst, "a", a, 1e-15, ai, 2.0, 1e-17, 1e-17, 1e-17)

	// 2 x 2 matrix
	io.Pf("\n----------------------------------(2 x 2)-----------------------------------\n")
	b := NewMatrixDeep2([][]float64{
		{1.0, 2.0},
		{3.0, 2.0},
	})
	bi := NewMatrixDeep2([][]float64{
		{-0.5, 0.5},
		{0.75, -0.25},
	})
	checkAiSmall(tst, "b", b, 1e-15, bi, -4.0, 1e-17, 1e-17, 1e-17)

	// 3 x 3 matrix
	io.Pf("\n----------------------------------(3 x 3)-----------------------------------\n")
	c := NewMatrixDeep2([][]float64{
		{+9.0, +3.0, 5.0},
		{-6.0, -9.0, 7.0},
		{-1.0, -8.0, 1.0},
	})
	ci := NewMatrixDeep2([][]float64{
		{+7.642276422764227e-02, -6.991869918699187e-02, +1.073170731707317e-01},
		{-1.626016260162601e-03, +2.276422764227642e-02, -1.512195121951219e-01},
		{+6.341463414634146e-02, +1.121951219512195e-01, -1.024390243902439e-01},
	})
	checkAiSmall(tst, "c", c, 1e-15, ci, 615.0, 1e-17, 1e-16, 1e-15)

	// 3 x 3 matrix
	io.Pf("\n----------------------------------(3 x 3)-----------------------------------\n")
	d := NewMatrixDeep2([][]float64{
		{1, 2, 3},
		{0, 4, 5},
		{1, 0, 6},
	})
	di := NewMatrixDeep2([][]float64{
		{12.0 / 11.0, -6.0 / 11.0, -1.0 / 11.0},
		{5.0 / 22.0, 3.0 / 22.0, -5.0 / 22.0},
		{-2.0 / 11.0, 1.0 / 11.0, 2.0 / 11.0},
	})
	checkAiSmall(tst, "d", d, 1e-15, di, 22.0, 1e-17, 1e-16, 1e-15)

	// 2 x 2 matrix
	io.Pf("\n----------------------------------(2 x 2)-----------------------------------\n")
	e := NewMatrixDeep2([][]float64{
		{1, 2},
		{3, 4},
	})
	ei := NewMatrixDeep2([][]float64{
		{-2.0, +1.0},
		{+1.5, -0.5},
	})
	checkAiSmall(tst, "e", e, 1e-15, ei, -2.0, 1e-17, 1e-17, 1e-17)

	// 3 x 3 matrix
	io.Pf("\n----------------------------------(3 x 3)-----------------------------------\n")
	f := NewMatrixDeep2([][]float64{
		{10, +1, +2},
		{+3, 20, +4},
		{+5, +6, 30},
	})
	fi := NewMatrixDeep2([][]float64{
		{+1.0423452768729642e-01, -3.2573289902280136e-03, -6.5146579804560255e-03},
		{-1.2667390517553386e-02, +5.2479189287006879e-02, -6.1527325370973572e-03},
		{-1.4838943177705389e-02, -9.9529496923633724e-03, +3.5649656170828804e-02},
	})
	checkAiSmall(tst, "f", f, 1e-15, fi, 5526.0, 1e-17, 1e-17, 1e-16)
}

func TestMatInv02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MatInv02. using OpenBLAS and larger matrices")

	// 1 x 1 matrix
	io.Pf("\n----------------------------------(1 x 1)-----------------------------------\n")
	A := NewMatrixDeep2([][]float64{
		{2.0},
	})
	Ai := NewMatrixDeep2([][]float64{
		{0.5},
	})
	checkAi(tst, "A", A, Ai, 2.0, 1e-17, 1e-17, 1e-17)

	// 2 x 2 matrix
	io.Pf("\n----------------------------------(2 x 2)-----------------------------------\n")
	B := NewMatrixDeep2([][]float64{
		{1.0, 2.0},
		{3.0, 2.0},
	})
	Bi := NewMatrixDeep2([][]float64{
		{-0.5, 0.5},
		{0.75, -0.25},
	})
	checkAi(tst, "B", B, Bi, -4.0, 1e-17, 1e-15, 1e-15)

	// 3 x 3 matrix
	io.Pf("\n----------------------------------(3 x 3)-----------------------------------\n")
	C := NewMatrixDeep2([][]float64{
		{+9.0, +3.0, 5.0},
		{-6.0, -9.0, 7.0},
		{-1.0, -8.0, 1.0},
	})
	Ci := NewMatrixDeep2([][]float64{
		{+7.642276422764227e-02, -6.991869918699187e-02, +1.073170731707317e-01},
		{-1.626016260162601e-03, +2.276422764227642e-02, -1.512195121951219e-01},
		{+6.341463414634146e-02, +1.121951219512195e-01, -1.024390243902439e-01},
	})
	checkAi(tst, "C", C, Ci, 615.0, 1.14e-13, 1e-16, 1e-15)

	// 4 x 4 matrix
	io.Pf("\n----------------------------------(4 x 4)-----------------------------------\n")
	D := NewMatrixDeep2([][]float64{
		{+3, +0, +2, -1},
		{+1, +2, +0, -2},
		{+4, +0, +6, -3},
		{+5, +0, +2, +0},
	})
	Di := NewMatrixDeep2([][]float64{
		{+0.6, 0.0, -0.2, 0.0},
		{-2.5, 0.5, +0.5, 1.0},
		{-1.5, 0.0, +0.5, 0.5},
		{-2.2, 0.0, +0.4, 1.0},
	})
	checkAi(tst, "D", D, Di, 20.0, 1e-14, 1e-15, 1e-15)

	// 5 x 5 matrix
	io.Pf("\n----------------------------------(5 x 5)-----------------------------------\n")
	a := NewMatrixDeep2([][]float64{
		{12, 28, 22, 20, +8},
		{+0, +3, +5, 17, 28},
		{56, +0, 23, +1, +0},
		{12, 29, 27, 10, +1},
		{+9, +4, 13, +8, 22},
	})
	ai := NewMatrixDeep2([][]float64{
		{+6.9128803717996279e-01, -7.4226114383340802e-01, -9.8756287260606410e-02, -6.9062496266472417e-01, +7.2471057693456553e-01},
		{+1.5936129795342968e+00, -1.7482347881148397e+00, -2.8304321334273236e-01, -1.5600769405383470e+00, +1.7164430532490673e+00},
		{-1.6345384165063759e+00, +1.7495848317224429e+00, +2.7469205863729274e-01, +1.6325730875377857e+00, -1.7065745928961444e+00},
		{-1.1177465024312745e+00, +1.3261729250546601e+00, +2.1243473793622566e-01, +1.1258168958554866e+00, -1.3325766717243535e+00},
		{+7.9976941733073770e-01, -8.9457712572131853e-01, -1.4770432850264653e-01, -8.0791149448632715e-01, +9.2990525800169743e-01},
	})
	detA := -167402.0
	checkAi(tst, "a", a, ai, detA, 1e-8, 1e-14, 1e-13)

	// 6 x 6 matrix
	io.Pf("\n----------------------------------(6 x 6)-----------------------------------\n")
	b := NewMatrixDeep2([][]float64{
		{+3.46540497998689445e-05, -1.39368151175265866e-05, -1.39368151175265866e-05, +0.00000000000000000e+00, 7.15957288480514429e-23, -2.93617909908697186e+02},
		{-1.39368151175265866e-05, +3.46540497998689445e-05, -1.39368151175265866e-05, +0.00000000000000000e+00, 7.15957288480514429e-23, -2.93617909908697186e+02},
		{-1.39368151175265866e-05, -1.39368151175265866e-05, +3.46540497998689445e-05, +0.00000000000000000e+00, 7.15957288480514429e-23, -2.93617909908697186e+02},
		{+0.00000000000000000e+00, +0.00000000000000000e+00, +0.00000000000000000e+00, +4.85908649173955311e-05, 0.00000000000000000e+00, +0.00000000000000000e+00},
		{+3.13760264822604860e-18, +3.13760264822604860e-18, +3.13760264822604860e-18, +0.00000000000000000e+00, 1.00000000000000000e+00, -1.93012141894243434e+07},
		{+0.00000000000000000e+00, +0.00000000000000000e+00, +0.00000000000000000e+00, -0.00000000000000000e+00, 0.00000000000000000e+00, +1.00000000000000000e+00},
	})
	bi := NewMatrixDeep2([][]float64{
		{+6.28811662297464645e+04, +4.23011662297464645e+04, +4.23011662297464645e+04, 0.00000000000000000e+00, -1.05591885817167332e-17, 4.33037966311565489e+07},
		{+4.23011662297464645e+04, +6.28811662297464645e+04, +4.23011662297464645e+04, 0.00000000000000000e+00, -1.05591885817167332e-17, 4.33037966311565489e+07},
		{+4.23011662297464645e+04, +4.23011662297464645e+04, +6.28811662297464645e+04, 0.00000000000000000e+00, -1.05591885817167348e-17, 4.33037966311565489e+07},
		{+0.00000000000000000e+00, +0.00000000000000000e+00, +0.00000000000000000e+00, 2.05800000000000000e+04, +0.00000000000000000e+00, 0.00000000000000000e+00},
		{-4.62744616057000471e-13, -4.62744616057000471e-13, -4.62744616057000471e-13, 0.00000000000000000e+00, +1.00000000000000000e+00, 1.93012141894243434e+07},
		{+0.00000000000000000e+00, +0.00000000000000000e+00, +0.00000000000000000e+00, 0.00000000000000000e+00, +0.00000000000000000e+00, 1.00000000000000000e+00},
	})
	detB := 0.0
	checkAi(tst, "b", b, bi, detB, 1e-15, 1e-8, 1e-9)
}

func TestMatSvd01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MatSvd01. SVD decomposition")

	// 2 x 2 matrix
	io.Pf("\n----------------------------------(2 x 2)-----------------------------------\n")

	a := NewMatrixDeep2([][]float64{
		{1, 2},
		{3, 4},
	})
	sA := []float64{5.4649857042190426e+00, 3.6596619062625751e-01}
	uA := NewMatrixDeep2([][]float64{
		{-4.0455358483375697e-01, -9.1451429567730447e-01},
		{-9.1451429567730458e-01, +4.0455358483375692e-01},
	})
	vtA := NewMatrixDeep2([][]float64{
		{-5.7604843676632078e-01, -8.1741556047036323e-01},
		{+8.1741556047036323e-01, -5.7604843676632078e-01},
	})
	checkSvd(tst, "a", a, sA, uA, vtA, 1e-14, 1e-15, 1e-15, 1e-14)

	// 3 x 3
	io.Pf("\n----------------------------------(3 x 3)-----------------------------------\n")
	b := NewMatrixDeep2([][]float64{
		{10, +1, +2},
		{+3, 20, +4},
		{+5, +6, 30},
	})
	sB := []float64{3.2864331638810626e+01, 1.7987569607075308e+01, 9.3478898990569999e+00}
	uB := NewMatrixDeep2([][]float64{
		{1.3141755556410686e-01, 4.4742200224954710e-02, -9.9031689958749303e-01},
		{3.6969554284917222e-01, 9.2470200356647714e-01, 9.0837273173509273e-02},
		{9.1981228067851462e-01, -3.7805335618027919e-01, 1.0498108493350734e-01},
	})
	vtB := NewMatrixDeep2([][]float64{
		{+2.1367614180500336e-01, 3.9691061543495060e-01, +8.9263904786782489e-01},
		{+7.4010067014497569e-02, 9.0453921735335097e-01, -4.1991822328912670e-01},
		{-9.7409702617544103e-01, 1.5579078157848747e-01, +1.6390306882827524e-01},
	})
	checkSvd(tst, "b", b, sB, uB, vtB, 1e-13, 1e-15, 1e-15, 1e-13)

	// 5 x 5
	io.Pf("\n----------------------------------(5 x 5)-----------------------------------\n")
	c := NewMatrixDeep2([][]float64{
		{12, 28, 22, 20, +8},
		{+0, +3, +5, 17, 28},
		{56, +0, 23, +1, +0},
		{12, 29, 27, 10, +1},
		{+9, +4, 13, +8, 22},
	})
	sC := []float64{7.6986806318205680e+01, 4.6904429440544916e+01, 3.2931871778592146e+01, 8.1528007049378086e+00, 1.7266616332203916e-01}
	uC := NewMatrixDeep2([][]float64{
		{-4.9131480299834873e-01, -3.9682713933839858e-01, +2.6940884231597306e-01, +5.5024083870837626e-01, +4.7517563167598015e-01},
		{-1.8436234721034561e-01, -4.5587438252398499e-01, -6.3940295404484160e-01, +2.7258142302445876e-01, -5.2445429016279244e-01},
		{-6.4666260791872432e-01, +7.1959191011975421e-01, -2.0909148108726730e-01, +1.1570356718901351e-01, -8.3116734093853104e-02},
		{-4.7937866776415455e-01, -2.6140673166893563e-01, +5.1576653173706322e-01, -4.6093491312241469e-01, -4.7263781496573321e-01},
		{-2.7684626365223813e-01, -2.2036508882136277e-01, -4.5699931653926035e-01, -6.3014766336473327e-01, +5.1851800447926977e-01},
	})
	vtC := NewMatrixDeep2([][]float64{
		{-6.5404770601013151e-01, -3.8083496832132990e-01, -5.6043632545196509e-01, -2.6778192710204507e-01, -2.0344603657478880e-01},
		{+6.4844738712320205e-01, -4.4646185600181121e-01, -9.3417253067825071e-02, -4.1240888414037769e-01, -4.4875374803990237e-01},
		{-1.9434098201682845e-01, +5.6949301281812137e-01, +1.7932721763630147e-01, -1.2720485785435223e-01, -7.6783459430605894e-01},
		{+2.3056394573648911e-01, +4.1309070476231217e-02, -5.5291118097154834e-01, +7.4868432524463224e-01, -2.8088224349796215e-01},
		{+2.4659568940095572e-01, +5.7411480305496243e-01, -5.8250778270734838e-01, -4.2603105176835093e-01, +2.9793453122548685e-01},
	})
	checkSvd(tst, "c", c, sC, uC, vtC, 1e-13, 1e-14, 1e-14, 1e-13)
}

func TestMatPseudo01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("MatPseudo01. Non-square matrices and pseudo inverse")

	// 4 x 3
	io.Pf("\n----------------------------------(4 x 3)-----------------------------------\n")
	a := NewMatrixDeep2([][]float64{
		{-5.773502691896260e-01, -5.773502691896260e-01, 1.000000000000000e+00},
		{+5.773502691896260e-01, -5.773502691896260e-01, 1.000000000000000e+00},
		{-5.773502691896260e-01, +5.773502691896260e-01, 1.000000000000000e+00},
		{+5.773502691896260e-01, +5.773502691896260e-01, 1.000000000000000e+00},
	})
	ai := NewMatrixDeep2([][]float64{
		{-4.330127018922192e-01, +4.330127018922192e-01, -4.330127018922192e-01, 4.330127018922192e-01},
		{-4.330127018922192e-01, -4.330127018922192e-01, +4.330127018922192e-01, 4.330127018922192e-01},
		{+2.500000000000000e-01, +2.500000000000000e-01, +2.500000000000000e-01, 2.500000000000000e-01},
	})
	sA := []float64{2, 1.15470053837925191e+00, 1.15470053837925191e+00}
	uA := NewMatrixDeep2([][]float64{
		{-0.5, -0.5, -0.5, +0.5},
		{-0.5, +0.5, -0.5, -0.5},
		{-0.5, -0.5, +0.5, -0.5},
		{-0.5, +0.5, +0.5, +0.5},
	})
	vtA := NewMatrixDeep2([][]float64{
		{0, 0, -1},
		{1, 0, 0},
		{0, 1, 0},
	})
	detA := 0.0
	checkAi(tst, "a", a, ai, detA, 1e-17, 1e-15, 1e-15)
	checkSvd(tst, "a", a, sA, uA, vtA, 1e-15, 1e-15, 1e-17, 1e-15)

	// 4 x 5
	io.Pf("\n----------------------------------(4 x 5)-----------------------------------\n")
	b := NewMatrixDeep2([][]float64{
		{1, 0, 0, 0, 2},
		{0, 0, 3, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 4, 0, 0, 0},
	})
	bi := NewMatrixDeep2([][]float64{
		{0.2, 0.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 1.0 / 4.0},
		{0.0, 1.0 / 3.0, 0.0, 0.0},
		{0.0, 0.0, 0.0, 0.0},
		{0.4, 0.0, 0.0, 0.0},
	})
	sB := []float64{4, 3, math.Sqrt(5.0), 0}
	detB := 0.0
	checkAi(tst, "b", b, bi, detB, 1e-17, 1e-15, 1e-15)
	checkSvd(tst, "b", b, sB, nil, nil, 1e-17, 1e-17, 1e-15, 1e-15)

	// 5 x 6
	io.Pf("\n----------------------------------(5 x 6)-----------------------------------\n")
	c := NewMatrixDeep2([][]float64{
		{12, 28, 22, 20, +8, 1},
		{+0, +3, +5, 17, 28, 1},
		{56, +0, 23, +1, +0, 1},
		{12, 29, 27, 10, +1, 1},
		{+9, +4, 13, +8, 22, 1},
	})
	ci := NewMatrixDeep2([][]float64{
		{+5.6387724512344639e-01, -6.0176177188969326e-01, -7.6500652148749224e-02, -5.6389938864086908e-01, +5.8595836573334192e-01},
		{+1.2836912791395787e+00, -1.4064756360496755e+00, -2.2890726327210095e-01, -1.2518220058421685e+00, +1.3789338004227019e+00},
		{-1.2866745075158739e+00, +1.3659857664770796e+00, +2.1392850711928030e-01, +1.2865799982753852e+00, -1.3277457214130808e+00},
		{-8.8185982449865485e-01, +1.0660542211012198e+00, +1.7123094548599221e-01, +8.9119882164767850e-01, -1.0756926383722674e+00},
		{+6.6698814093525072e-01, -7.4815557352521045e-01, -1.2451059750508876e-01, -6.7584431870600359e-01, +7.8530451101142418e-01},
		{-1.1017522295492406e+00, +1.2149323757487696e+00, +1.9244991110051662e-01, +1.0958269819071325e+00, -1.1998242501940171e+00},
	})
	detC := 0.0
	checkAi(tst, "c", c, ci, detC, 1e-17, 1e-13, 1e-12)
	checkSvd(tst, "c", c, nil, nil, nil, 1e-17, 1e-17, 1e-17, 1e-13)

	// 8 x 6
	io.Pf("\n----------------------------------(8 x 6)-----------------------------------\n")
	d := NewMatrixDeep2([][]float64{
		{64, +2, +3, 61, 60, +6},
		{+9, 55, 54, 12, 13, 51},
		{17, 47, 46, 20, 21, 43},
		{40, 26, 27, 37, 36, 30},
		{32, 34, 35, 29, 28, 38},
		{41, 23, 22, 44, 45, 19},
		{49, 15, 14, 52, 53, 11},
		{+8, 58, 59, +5, +4, 62},
	})
	sD := []float64{2.25169577993700130e+02, 1.27186528905283367e+02, 1.17578914421132179e+01, 1.81235447053960281e-14, 9.59676789459164647e-15, 5.90626950718289933e-15}
	detD := 0.0
	checkAi(tst, "d", d, nil, detD, 1e-17, 1e-17, 1e-13)
	checkSvd(tst, "d", d, sD, nil, nil, 1e-13, 1e-17, 1e-17, 1e-12)
}

func TestCondNum01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("CondNum01. condition number of matrix using inverse")

	a := NewMatrixDeep2([][]float64{
		{1, 2},
		{2, 3.999},
	})

	cIa := MatCondNum(a, "I")
	chk.Float64(tst, "condI(a) ", 1e-8, cIa, 35988.001)

	b := NewMatrixDeep2([][]float64{
		{1, 2},
		{2, 3},
	})

	cIb := MatCondNum(b, "I")
	cFb := MatCondNum(b, "F")
	chk.Float64(tst, "condI(b) ", 1e-17, cIb, 25.0)
	chk.Float64(tst, "condF(b) ", 1e-14, cFb, 18.0)
}
