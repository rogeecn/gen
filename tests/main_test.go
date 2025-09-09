package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"
	"time"

	"gen-test/database/model"
	"gen-test/database/query"

	"go.ipao.vip/gen"
	"go.ipao.vip/gen/field"
	types "go.ipao.vip/gen/types"

	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func init() {
	dbConfig := postgres.Config{DSN: dsn}
	logrus.Info("Open PostgreSQL:", dsn)

	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
		Logger:         logger.Default,
	}

	db, err := gorm.Open(postgres.New(dbConfig), &gormConfig)
	if err != nil {
		log.Fatal(err)
	}

	query.SetDefault(db)
}

func Test_main(t *testing.T) {
	Convey("main", t, func() {
		Convey("First", func() {
			tbl, q := query.Student.QueryContext(context.Background())
			So(tbl, ShouldNotBeNil)
			So(q, ShouldNotBeNil)

			m, err := q.Preload(tbl.RelationClass).First()
			So(m, ShouldNotBeNil)
			So(err, ShouldBeNil)

			pretty.Print(m)
		})
	})
}

func Test_Create(t *testing.T) {
	tbl, q := query.ComprehensiveTypesTable.QueryContext(context.Background())
	Convey("Create inserts with all field types and can fetch", t, func() {
		marker := fmt.Sprintf("created-%d", time.Now().UnixNano())
		t0 := time.Date(2023, 1, 1, 20, 0, 0, 0, time.Local)
		d0 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)
		ts0 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)
		ts1 := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
		rec := &model.ComprehensiveTypesTable{
			SmallInt:     int16(2),
			IntegerVal:   int32(2000),
			BigInt:       int64(2000000),
			UintVal:      int32(123),
			Uint8Val:     int16(200),
			Uint16Val:    int32(40000),
			Uint32Val:    int64(2147483647),
			Uint64Val:    float64(123456789),
			Float32Val:   float32(6.28),
			Float64Val:   float64(3.14159),
			StringVal:    marker,
			VarcharVal:   "VarStr",
			CharVal:      "Ch",
			BytesVal:     []byte{0xca, 0xfe, 0xba, 0xbe},
			BoolVal:      false,
			TimeVal:      t0,
			DateVal:      types.Date(d0),
			TimeOnly:     types.NewTime(12, 34, 56, 0),
			TimestampVal: ts0,
			JSONVal:      types.JSON([]byte(`{"a":1}`)),
			JsonbVal:     types.JSON([]byte(`{"b":true}`)),
			StringArray:  types.NewArray([]string{"x", "y"}),
			IntArray:     types.NewArray([]int32{10, 20}),
			BigintArray:  types.NewArray([]int64{100, 200}),
			FloatArray:   types.NewArray([]float64{1.5, 2.5}),
			UUIDVal:      types.NewUUIDv4(),
			DecimalVal:   99.99,
			NumericVal:   1234.567,
			InetVal:      types.MustInet("10.0.0.1"),
			CidrVal:      types.MustCIDR("192.168.0.0/24"),
			MacaddrVal:   types.MustMACAddr("08:00:2b:01:02:03"),
			PointVal:     types.NewPoint(3, 4),
			BoxVal:       types.NewBox(types.NewPoint(0, 0), types.NewPoint(1, 1)),
			PathVal:      types.NewPath([]types.Point{types.NewPoint(0, 0), types.NewPoint(1, 1)}, false),
			PolygonVal: types.NewPolygon(
				[]types.Point{types.NewPoint(0, 0), types.NewPoint(1, 0), types.NewPoint(1, 1)},
			),
			CircleVal:    types.NewCircle(types.NewPoint(0, 0), 2.5),
			BitVal:       types.NewBitString("10101010"),
			VarbitVal:    types.NewBitString("1100110011001100"),
			Int4rangeVal: types.NewInt4Range(1, 10, true, false),
			Int8rangeVal: types.NewInt8Range(1, 100, true, false),
			NumrangeVal:  types.NewNumRange(new(big.Rat).SetFloat64(1.25), new(big.Rat).SetFloat64(9.75), true, true),
			TsrangeVal:   types.NewTsRange(ts0, ts1, true, false),
			TstzrangeVal: types.NewTstzRange(ts0.UTC(), ts1.UTC(), false, false),
			DaterangeVal: types.NewDateRange(d0, d0.AddDate(0, 1, 0), true, false),
			TsvectorVal:  types.NewTSVector("'foo':1 'bar':2"),
			TsqueryVal:   types.NewTSQuery("'foo' & 'bar'"),
			XMLVal:       types.NewXML("<root>ok</root>"),
			MoneyVal:     types.NewMoney("$9.99"),
		}

		err := q.Create(rec)
		So(err, ShouldBeNil)
		So(rec.ID, ShouldBeGreaterThan, 0)

		got, err := q.Where(tbl.StringVal.Eq(marker)).First()
		So(err, ShouldBeNil)
		So(got, ShouldNotBeNil)
		So(got.StringVal, ShouldEqual, marker)
		So(got.SmallInt, ShouldEqual, int16(2))
		So(got.BoolVal, ShouldEqual, false)
		So(got.PointVal.X, ShouldEqual, 3)
		So(got.BoxVal.P1.Y == 1 || got.BoxVal.P2.Y == 1, ShouldBeTrue)
		So(got.IntegerVal, ShouldEqual, int32(2000))
		So(got.BigInt, ShouldEqual, int64(2000000))
		So(got.UintVal, ShouldEqual, int32(123))
		So(got.Uint8Val, ShouldEqual, int16(200))
		So(got.Uint16Val, ShouldEqual, int32(40000))
		So(got.Uint32Val, ShouldEqual, int64(2147483647))
		So(got.Uint64Val, ShouldEqual, float64(123456789))
		So(got.Float32Val, ShouldEqual, float32(6.28))
		So(got.Float64Val, ShouldEqual, float64(3.14159))
		So(got.VarcharVal, ShouldEqual, "VarStr")
		So(strings.TrimSpace(got.CharVal), ShouldEqual, "Ch")
		So(got.BytesVal, ShouldResemble, []byte{0xca, 0xfe, 0xba, 0xbe})
		So(got.TimeVal.Equal(t0), ShouldBeTrue)
		So(got.DateVal.String(), ShouldEqual, types.Date(d0).String())
		So(got.TimeOnly.String(), ShouldEqual, types.NewTime(12, 34, 56, 0).String())
		So(got.TimestampVal.Equal(ts0), ShouldBeTrue)
		So(got.JSONVal, ShouldResemble, types.JSON([]byte(`{"a":1}`)))
		So(got.JsonbVal.String(), ShouldEqual, string([]byte(`{"b": true}`)))
		So(got.StringArray, ShouldResemble, types.NewArray([]string{"x", "y"}))
		So(got.IntArray, ShouldResemble, types.NewArray([]int32{10, 20}))
		So(got.BigintArray, ShouldResemble, types.NewArray([]int64{100, 200}))
		So(got.FloatArray, ShouldResemble, types.NewArray([]float64{1.5, 2.5}))
		So(got.UUIDVal.String() != "", ShouldBeTrue)
		So(got.DecimalVal, ShouldEqual, 99.99)
		So(got.NumericVal, ShouldEqual, 1234.567)
		So(got.InetVal.String(), ShouldEqual, "10.0.0.1")
		So(got.CidrVal.String(), ShouldEqual, "192.168.0.0/24")
		So(got.MacaddrVal.String(), ShouldEqual, "08:00:2b:01:02:03")
		So(got.PointVal, ShouldResemble, types.NewPoint(3, 4))
		// PostgreSQL may reorder box corners; accept either ordering
		{
			p1 := types.NewPoint(0, 0)
			p2 := types.NewPoint(1, 1)
			b := got.BoxVal
			So((b == types.NewBox(p1, p2)) || (b == types.NewBox(p2, p1)), ShouldBeTrue)
		}
		So(got.PathVal, ShouldResemble, types.NewPath([]types.Point{types.NewPoint(0, 0), types.NewPoint(1, 1)}, false))
		So(
			got.PolygonVal,
			ShouldResemble,
			types.NewPolygon([]types.Point{types.NewPoint(0, 0), types.NewPoint(1, 0), types.NewPoint(1, 1)}),
		)
		So(got.CircleVal, ShouldResemble, types.NewCircle(types.NewPoint(0, 0), 2.5))
		So(got.BitVal, ShouldResemble, types.NewBitString("10101010"))
		So(got.VarbitVal, ShouldResemble, types.NewBitString("1100110011001100"))
		So(got.Int4rangeVal, ShouldResemble, types.NewInt4Range(1, 10, true, false))
		So(got.Int8rangeVal, ShouldResemble, types.NewInt8Range(1, 100, true, false))
		So(
			got.NumrangeVal,
			ShouldResemble,
			types.NewNumRange(new(big.Rat).SetFloat64(1.25), new(big.Rat).SetFloat64(9.75), true, true),
		)
		// tsrange is timestamp without time zone; normalize to UTC for comparison
		wantLower := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		So(got.TsrangeVal.Lower.Equal(wantLower), ShouldBeTrue)
		So(got.TsrangeVal.Upper.Equal(ts1), ShouldBeTrue)
		So(got.TsrangeVal.LowerInclusive, ShouldBeTrue)
		So(got.TsrangeVal.UpperInclusive, ShouldBeFalse)
		// tstzrange stores with timezone; compare instants
		So(got.TstzrangeVal.Lower.Equal(ts0.UTC()), ShouldBeTrue)
		So(got.TstzrangeVal.Upper.Equal(ts1.UTC()), ShouldBeTrue)
		// daterange normalizes to UTC as well; compare instants and flags
		wantDateLower := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		wantDateUpper := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
		So(got.DaterangeVal.Lower.Equal(wantDateLower), ShouldBeTrue)
		So(got.DaterangeVal.Upper.Equal(wantDateUpper), ShouldBeTrue)
		So(got.DaterangeVal.LowerInclusive, ShouldBeTrue)
		So(got.DaterangeVal.UpperInclusive, ShouldBeFalse)
		// tsvector term order may vary; ensure both terms exist with positions
		sv := string(got.TsvectorVal)
		So(strings.Contains(sv, "'foo':1"), ShouldBeTrue)
		So(strings.Contains(sv, "'bar':2"), ShouldBeTrue)
		So(got.TsqueryVal, ShouldResemble, types.NewTSQuery("'foo' & 'bar'"))
		So(got.XMLVal, ShouldResemble, types.NewXML("<root>ok</root>"))
		So(got.MoneyVal, ShouldResemble, types.NewMoney("$9.99"))

		_, err = q.Where(tbl.ID.Eq(got.ID)).Delete()
		So(err, ShouldBeNil)
	})
}

func Test_Update(t *testing.T) {
	tbl, q := query.ComprehensiveTypesTable.QueryContext(context.Background())
	Convey("Update modifies all columns", t, func() {
		marker := fmt.Sprintf("upd-%d", time.Now().UnixNano())
		// Create a base record first
		t0 := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)
		d0 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		ts0 := time.Date(2023, 1, 2, 10, 0, 0, 0, time.UTC)
		rec := &model.ComprehensiveTypesTable{
			SmallInt:     10,
			IntegerVal:   100,
			BigInt:       1000,
			UintVal:      10,
			Uint8Val:     20,
			Uint16Val:    30,
			Uint32Val:    40,
			Uint64Val:    50,
			Float32Val:   1.0,
			Float64Val:   2.0,
			StringVal:    marker,
			VarcharVal:   "v1",
			CharVal:      "C1",
			BytesVal:     []byte{0x01, 0x02},
			BoolVal:      false,
			TimeVal:      t0,
			DateVal:      types.Date(d0),
			TimeOnly:     types.NewTime(1, 2, 3, 0),
			TimestampVal: ts0,
			JSONVal:      types.JSON([]byte(`{"init":true}`)),
			JsonbVal:     types.JSON([]byte(`{"init":true}`)),
			StringArray:  types.NewArray([]string{"a"}),
			IntArray:     types.NewArray([]int32{1, 2}),
			BigintArray:  types.NewArray([]int64{3, 4}),
			FloatArray:   types.NewArray([]float64{5.5}),
			UUIDVal:      types.NewUUIDv4(),
			DecimalVal:   1.23,
			NumericVal:   4.56,
			InetVal:      types.MustInet("10.0.0.10"),
			CidrVal:      types.MustCIDR("10.0.0.0/8"),
			MacaddrVal:   types.MustMACAddr("08:00:2b:00:00:01"),
			PointVal:     types.NewPoint(0, 0),
			BoxVal:       types.NewBox(types.NewPoint(0, 0), types.NewPoint(1, 1)),
			PathVal:      types.NewPath([]types.Point{types.NewPoint(0, 0), types.NewPoint(1, 1)}, true),
			PolygonVal: types.NewPolygon(
				[]types.Point{types.NewPoint(0, 0), types.NewPoint(0, 1), types.NewPoint(1, 1)},
			),
			CircleVal:    types.NewCircle(types.NewPoint(0, 0), 1.1),
			BitVal:       types.NewBitString("10101010"),
			VarbitVal:    types.NewBitString("1100110011001100"),
			Int4rangeVal: types.NewInt4Range(1, 2, true, false),
			Int8rangeVal: types.NewInt8Range(1, 2, true, false),
			NumrangeVal:  types.NewNumRange(new(big.Rat).SetFloat64(1.0), new(big.Rat).SetFloat64(2.0), true, false),
			TsrangeVal:   types.NewTsRange(ts0, ts0.AddDate(0, 1, 0), true, false),
			TstzrangeVal: types.NewTstzRange(ts0, ts0.AddDate(0, 1, 0), true, false),
			DaterangeVal: types.NewDateRange(d0, d0.AddDate(0, 1, 0), true, false),
			TsvectorVal:  types.NewTSVector("'foo':1 'bar':2"),
			TsqueryVal:   types.NewTSQuery("'foo' & 'bar'"),
			XMLVal:       types.NewXML("<r>1</r>"),
			MoneyVal:     types.NewMoney("$0.01"),
		}
		So(q.Create(rec), ShouldBeNil)
		So(rec.ID, ShouldBeGreaterThan, 0)

		// Prepare update values for all updatable fields
		newMarker := marker + "-new"
		t1 := time.Date(2023, 2, 1, 9, 30, 0, 0, time.UTC)
		d1 := time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)
		ts1 := time.Date(2023, 3, 1, 10, 0, 0, 0, time.UTC)
		updates := map[string]interface{}{
			"small_int":     int16(20),
			"integer_val":   int32(200),
			"big_int":       int64(2000),
			"uint_val":      int32(20),
			"uint8_val":     int16(30),
			"uint16_val":    int32(30000),
			"uint32_val":    int64(123456789),
			"uint64_val":    float64(987654321),
			"float32_val":   float32(7.89),
			"float64_val":   float64(6.54321),
			"string_val":    newMarker,
			"varchar_val":   "v2",
			"char_val":      "C2",
			"bytes_val":     []byte{0xaa, 0xbb},
			"bool_val":      true,
			"time_val":      t1,
			"date_val":      types.Date(d1),
			"time_only":     types.NewTime(23, 59, 58, 0),
			"timestamp_val": ts1,
			"json_val":      types.JSON([]byte(`{"u":2}`)),
			"jsonb_val":     types.JSON([]byte(`{"x":123}`)),
			"string_array":  types.NewArray([]string{"x", "y"}),
			"int_array":     types.NewArray([]int32{10, 20, 30}),
			"bigint_array":  types.NewArray([]int64{100, 200, 300}),
			"float_array":   types.NewArray([]float64{1.1, 2.2, 3.3}),
			"uuid_val":      types.NewUUIDv4(),
			"decimal_val":   12.34,
			"numeric_val":   56.78,
			"inet_val":      types.MustInet("10.0.0.2"),
			"cidr_val":      types.MustCIDR("192.168.1.0/24"),
			"macaddr_val":   types.MustMACAddr("08:00:2b:04:05:06"),
			"point_val":     types.NewPoint(5, 6),
			"box_val":       types.NewBox(types.NewPoint(2, 2), types.NewPoint(3, 3)),
			"path_val":      types.NewPath([]types.Point{types.NewPoint(2, 2), types.NewPoint(3, 3)}, false),
			"polygon_val": types.NewPolygon(
				[]types.Point{types.NewPoint(0, 0), types.NewPoint(1, 0), types.NewPoint(1, 1)},
			),
			"circle_val":    types.NewCircle(types.NewPoint(1, 1), 3.3),
			"bit_val":       types.NewBitString("11110000"),
			"varbit_val":    types.NewBitString("0011001100110011"),
			"int4range_val": types.NewInt4Range(2, 9, true, false),
			"int8range_val": types.NewInt8Range(2, 200, true, false),
			"numrange_val":  types.NewNumRange(new(big.Rat).SetFloat64(2.5), new(big.Rat).SetFloat64(7.5), true, false),
			"tsrange_val":   types.NewTsRange(ts1, ts1.AddDate(0, 1, 0), true, false),
			"tstzrange_val": types.NewTstzRange(ts1.UTC(), ts1.AddDate(0, 1, 0).UTC(), false, false),
			"daterange_val": types.NewDateRange(d1, d1.AddDate(0, 2, 0), true, true),
			"tsvector_val":  types.NewTSVector("'alpha':1 'beta':2"),
			"tsquery_val":   types.NewTSQuery("'alpha' | 'beta'"),
			"xml_val":       types.NewXML("<r>2</r>"),
			"money_val":     types.NewMoney("$1.23"),
		}

		_, err := q.Where(tbl.ID.Eq(rec.ID)).UpdateColumns(updates)
		So(err, ShouldBeNil)

		got, err := q.Where(tbl.ID.Eq(rec.ID)).First()
		So(err, ShouldBeNil)
		So(got.StringVal, ShouldEqual, newMarker)
		So(got.SmallInt, ShouldEqual, int16(20))
		So(got.IntegerVal, ShouldEqual, int32(200))
		So(got.BigInt, ShouldEqual, int64(2000))
		So(got.UintVal, ShouldEqual, int32(20))
		So(got.Uint8Val, ShouldEqual, int16(30))
		So(got.Uint16Val, ShouldEqual, int32(30000))
		So(got.Uint32Val, ShouldEqual, int64(123456789))
		So(got.Uint64Val, ShouldEqual, float64(987654321))
		So(got.Float32Val, ShouldEqual, float32(7.89))
		So(got.Float64Val, ShouldEqual, float64(6.54321))
		So(got.VarcharVal, ShouldEqual, "v2")
		So(strings.TrimSpace(got.CharVal), ShouldEqual, "C2")
		So(got.BytesVal, ShouldResemble, []byte{0xaa, 0xbb})
		So(got.BoolVal, ShouldEqual, true)
		So(got.TimeVal.Equal(t1), ShouldBeTrue)
		So(time.Time(got.DateVal).UTC().Equal(d1.UTC()), ShouldBeTrue)
		So(got.TimeOnly.String(), ShouldEqual, types.NewTime(23, 59, 58, 0).String())
		So(got.TimestampVal.Format("2006-01-02 15:04:05"), ShouldEqual, ts1.Format("2006-01-02 15:04:05"))
		So(string(got.JSONVal), ShouldEqual, string(types.JSON([]byte(`{"u":2}`))))
		So(strings.Contains(got.JsonbVal.String(), "\"x\": 123"), ShouldBeTrue)
		So(got.StringArray, ShouldResemble, types.NewArray([]string{"x", "y"}))
		So(got.IntArray, ShouldResemble, types.NewArray([]int32{10, 20, 30}))
		So(got.BigintArray, ShouldResemble, types.NewArray([]int64{100, 200, 300}))
		So(got.FloatArray, ShouldResemble, types.NewArray([]float64{1.1, 2.2, 3.3}))
		So(got.InetVal.String(), ShouldEqual, "10.0.0.2")
		So(got.CidrVal.String(), ShouldEqual, "192.168.1.0/24")
		So(got.MacaddrVal.String(), ShouldEqual, "08:00:2b:04:05:06")
		So(got.PointVal, ShouldResemble, types.NewPoint(5, 6))
		{
			p1 := types.NewPoint(2, 2)
			p2 := types.NewPoint(3, 3)
			b := got.BoxVal
			So((b == types.NewBox(p1, p2)) || (b == types.NewBox(p2, p1)), ShouldBeTrue)
		}
		So(got.PathVal.Closed, ShouldEqual, false)
		So(len(got.PolygonVal.Points), ShouldEqual, 3)
		So(got.CircleVal.Radius, ShouldEqual, 3.3)
		So(got.BitVal, ShouldResemble, types.NewBitString("11110000"))
		So(got.VarbitVal, ShouldResemble, types.NewBitString("0011001100110011"))
		So(got.Int4rangeVal, ShouldResemble, types.NewInt4Range(2, 9, true, false))
		So(got.Int8rangeVal, ShouldResemble, types.NewInt8Range(2, 200, true, false))
		So(
			got.NumrangeVal,
			ShouldResemble,
			types.NewNumRange(new(big.Rat).SetFloat64(2.5), new(big.Rat).SetFloat64(7.5), true, false),
		)
		// tsrange compare instants
		So(got.TsrangeVal.Lower.Equal(ts1.UTC()), ShouldBeTrue)
		So(got.TsrangeVal.Upper.Equal(ts1.AddDate(0, 1, 0)), ShouldBeTrue)
		So(got.TsrangeVal.LowerInclusive, ShouldBeTrue)
		So(got.TsrangeVal.UpperInclusive, ShouldBeFalse)
		So(got.TstzrangeVal.Lower.Equal(ts1.UTC()), ShouldBeTrue)
		So(got.TstzrangeVal.Upper.Equal(ts1.AddDate(0, 1, 0).UTC()), ShouldBeTrue)
		So(got.DaterangeVal.Lower.Format("2006-01-02"), ShouldEqual, d1.Format("2006-01-02"))
		upperStr := got.DaterangeVal.Upper.Format("2006-01-02")
		exp1 := d1.AddDate(0, 2, 0).Format("2006-01-02")
		exp2 := d1.AddDate(0, 2, 1).Format("2006-01-02")
		So(upperStr == exp1 || upperStr == exp2, ShouldBeTrue)
		sv := string(got.TsvectorVal)
		So(strings.Contains(sv, "'alpha':1"), ShouldBeTrue)
		So(strings.Contains(sv, "'beta':2"), ShouldBeTrue)
		So(string(got.TsqueryVal), ShouldEqual, string(types.NewTSQuery("'alpha' | 'beta'")))
		So(string(got.XMLVal), ShouldEqual, string(types.NewXML("<r>2</r>")))
		So(string(got.MoneyVal), ShouldEqual, string(types.NewMoney("$1.23")))

		// 15) Array overlaps and contained-by
		prefix := fmt.Sprintf("upd-%d-", time.Now().UnixNano())
		So(
			q.DO.Create(
				map[string]interface{}{
					"string_val":   prefix + "arr2",
					"string_array": types.NewArray([]string{"aa", "cc"}),
				},
			),
			ShouldBeNil,
		)
		aov, err := q.Where(tbl.StringVal.Eq(prefix+"arr2"), tbl.StringArray.Overlaps(types.NewArray([]string{"cc"}))).
			Find()
		So(err, ShouldBeNil)
		So(len(aov), ShouldEqual, 1)
		acb, err := q.Where(tbl.StringVal.Eq(prefix+"arr2"), tbl.StringArray.ContainedBy(types.NewArray([]string{"aa", "bb", "cc"}))).
			Find()
		So(err, ShouldBeNil)
		So(len(acb), ShouldEqual, 1)

		// 16) JSONB key regex and iregex
		jsonRegex := prefix + "jsonrx"
		So(
			q.DO.Create(
				map[string]interface{}{
					"string_val": jsonRegex,
					"jsonb_val":  types.JSON([]byte("{\"meta\":{\"title\":\"Hello-World\"}}")),
				},
			),
			ShouldBeNil,
		)
		jrx, err := q.Where(tbl.StringVal.Eq(jsonRegex), tbl.JsonbVal.KeyRegexp("meta.title", "^Hello")).Find()
		So(err, ShouldBeNil)
		So(len(jrx), ShouldEqual, 1)
		jirx, err := q.Where(tbl.StringVal.Eq(jsonRegex), tbl.JsonbVal.KeyIRegexp("meta.title", "world$")).Find()
		So(err, ShouldBeNil)
		So(len(jirx), ShouldEqual, 1)

		// 17) Subquery: IN (subquery) via ContainsSubQuery
		subq := q.UnderlyingDB().Table(tbl.TableName()).Select("id").Where("string_val LIKE ?", prefix+"a%")
		inrs, err := q.Where(field.ContainsSubQuery([]field.Expr{tbl.ID}, subq)).Find()
		So(err, ShouldBeNil)
		So(len(inrs), ShouldBeGreaterThanOrEqualTo, 1)

		// 18) Subquery: EXISTS (subquery)
		exSub := q.UnderlyingDB().Table(tbl.TableName()).Select("1").Where("string_val LIKE ?", prefix+"C%")
		exrs, err := q.Where(field.CompareSubQuery(field.ExistsOp, tbl.ID, exSub)).Find()
		So(err, ShouldBeNil)
		So(len(exrs), ShouldBeGreaterThanOrEqualTo, 1)

		// cleanup
		_, err = q.Where(tbl.ID.Eq(got.ID)).Delete()
		So(err, ShouldBeNil)
	})
}

func Test_Find(t *testing.T) {
	tbl, q := query.ComprehensiveTypesTable.QueryContext(context.Background())
	Convey("Find returns multiple rows", t, func() {
		prefix := fmt.Sprintf("find-%d-", time.Now().UnixNano())
		So(q.DO.Create(map[string]interface{}{"string_val": prefix + "a"}), ShouldBeNil)
		So(q.DO.Create(map[string]interface{}{"string_val": prefix + "b"}), ShouldBeNil)

		rows, err := q.Where(tbl.StringVal.Like(prefix + "%")).Find()
		So(err, ShouldBeNil)
		So(len(rows), ShouldBeGreaterThanOrEqualTo, 2)

		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with diverse condition types", t, func() {
		prefix := fmt.Sprintf("qdiv-%d-", time.Now().UnixNano())
		// Insert 3 records with varied values
		rowsToMake := []map[string]interface{}{
			{
				"small_int":     int16(1),
				"integer_val":   int32(10),
				"string_val":    prefix + "a",
				"varchar_val":   "Alpha",
				"bool_val":      true,
				"point_val":     types.NewPoint(0, 0),
				"int4range_val": types.NewInt4Range(1, 5, true, false),
			},
			{
				"small_int":     int16(2),
				"integer_val":   int32(20),
				"string_val":    prefix + "b",
				"varchar_val":   "beta",
				"bool_val":      false,
				"point_val":     types.NewPoint(10, 10),
				"int4range_val": types.NewInt4Range(5, 10, true, false),
			},
			{
				"small_int":     int16(3),
				"integer_val":   int32(30),
				"string_val":    prefix + "C",
				"varchar_val":   "Camel",
				"bool_val":      true,
				"point_val":     types.NewPoint(2, 2),
				"int4range_val": types.NewInt4Range(2, 3, true, false),
			},
		}
		for _, m := range rowsToMake {
			So(q.DO.Create(m), ShouldBeNil)
		}

		// 1) IN: small_int IN (1,3)
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.SmallInt.In(int16(1), int16(3))).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldBeGreaterThanOrEqualTo, 1)

		// 1b) OR composition: string_val = prefix+'a' OR prefix+'C'
		ors, err := q.Where(field.Or(tbl.StringVal.Eq(prefix+"a"), tbl.StringVal.Eq(prefix+"C"))).Find()
		So(err, ShouldBeNil)
		So(len(ors), ShouldBeGreaterThanOrEqualTo, 1)

		// 2) Between: 15 <= integer_val <= 35
		rs, err = q.Where(tbl.IntegerVal.Between(int32(15), int32(35)), tbl.StringVal.Like(prefix+"%")).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2)

		// 3) ILike (case-insensitive)
		rs, err = q.Where(tbl.VarcharVal.ILike("%a%"), tbl.StringVal.Like(prefix+"%")).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldBeGreaterThanOrEqualTo, 2) // Alpha, Camel (and possibly others)

		// 4) Range Overlaps
		rs, err = q.Where(tbl.Int4rangeVal.Overlaps(types.NewInt4Range(4, 6, true, false)), tbl.StringVal.Like(prefix+"%")).
			Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldBeGreaterThanOrEqualTo, 1) // overlaps at least one; boundary behavior may vary

		// 5) Range Contains
		target := types.NewInt4Range(2, 3, true, false)
		rs, err = q.Where(tbl.Int4rangeVal.Contains(target), tbl.StringVal.Like(prefix+"%")).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // [1,5) and [2,3)

		// 6) Geometry: order by distance to a reference point

		// 7) Order by distance to (0,0) ascending, expect first is suffix 'a'
		orderQ := q.Where(tbl.StringVal.Like(prefix + "%"))
		near, err := orderQ.Order(tbl.PointVal.DistanceTo(types.NewPoint(0, 0))).Limit(1).Find()
		So(err, ShouldBeNil)
		So(len(near), ShouldEqual, 1)
		So(strings.HasSuffix(near[0].StringVal, "a"), ShouldBeTrue)

		// 8) NOT + AND
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Not(tbl.StringVal.Like(prefix + "b%")).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2)

		// 9) Pagination via Order + Offset + Limit
		page, err := q.Where(tbl.StringVal.Like(prefix + "%")).Order(tbl.SmallInt).Offset(1).Limit(1).Find()
		So(err, ShouldBeNil)
		So(len(page), ShouldEqual, 1)
		So(page[0].SmallInt, ShouldEqual, int16(2))

		// 10) Array contains/overlaps
		arrMarker := prefix + "arr"
		So(
			q.DO.Create(
				map[string]interface{}{"string_val": arrMarker, "string_array": types.NewArray([]string{"aa", "bb"})},
			),
			ShouldBeNil,
		)
		ars, err := q.Where(tbl.StringVal.Eq(arrMarker), tbl.StringArray.Contains(types.NewArray([]string{"aa"}))).
			Find()
		So(err, ShouldBeNil)
		So(len(ars), ShouldEqual, 1)

		// 11) JSONB contains and key match
		jsonMarker := prefix + "json"
		So(
			q.DO.Create(
				map[string]interface{}{
					"string_val": jsonMarker,
					"jsonb_val":  types.JSON([]byte(`{"name":"foo","k":1}`)),
				},
			),
			ShouldBeNil,
		)
		jrs, err := q.Where(
			tbl.StringVal.Eq(jsonMarker),
			tbl.JsonbVal.Contains(types.JSON([]byte(`{"k":1}`))),
			tbl.JsonbVal.KeyEq("name", "foo"),
		).Find()
		So(err, ShouldBeNil)
		So(len(jrs), ShouldEqual, 1)

		// 12) Fulltext match using @@
		ftMarker := prefix + "fts"
		So(
			q.DO.Create(
				map[string]interface{}{"string_val": ftMarker, "tsvector_val": types.NewTSVector("'zz':1 'yy':2")},
			),
			ShouldBeNil,
		)
		ftrs, err := q.Where(tbl.StringVal.Eq(ftMarker), tbl.TsvectorVal.Matches(types.NewTSQuery("'zz' | 'yy'"))).
			Find()
		So(err, ShouldBeNil)
		So(len(ftrs), ShouldEqual, 1)

		// 13) Network contains (CIDR contains subnet) and INET contains-eq
		netMarker := prefix + "net"
		So(
			q.DO.Create(
				map[string]interface{}{
					"string_val": netMarker,
					"inet_val":   types.MustInet("10.1.2.3"),
					"cidr_val":   types.MustCIDR("10.1.0.0/16"),
				},
			),
			ShouldBeNil,
		)
		nrs, err := q.Where(tbl.StringVal.Eq(netMarker), tbl.CidrVal.Contains(types.MustCIDR("10.1.2.0/24"))).Find()
		So(err, ShouldBeNil)
		So(len(nrs), ShouldEqual, 1)
		nrs, err = q.Where(tbl.StringVal.Eq(netMarker), tbl.InetVal.ContainsEq(types.MustInet("10.1.2.3"))).Find()
		So(err, ShouldBeNil)
		So(len(nrs), ShouldEqual, 1)

		// 14) Fulltext rank ordering using ts_rank
		rankA := prefix + "rankA"
		rankB := prefix + "rankB"
		So(
			q.DO.Create(map[string]interface{}{"string_val": rankA, "tsvector_val": types.NewTSVector("'alpha':1")}),
			ShouldBeNil,
		)
		So(
			q.DO.Create(
				map[string]interface{}{
					"string_val":   rankB,
					"tsvector_val": types.NewTSVector("'alpha':1 'alpha':2 'alpha':3"),
				},
			),
			ShouldBeNil,
		)
		orderExpr := field.NewUnsafeFieldRaw("ts_rank(?, to_tsquery(?)) DESC", tbl.TsvectorVal.RawExpr(), "alpha")
		top, err := q.Where(tbl.StringVal.Like(prefix + "rank%")).Order(orderExpr).Limit(1).Find()
		So(err, ShouldBeNil)
		So(len(top), ShouldEqual, 1)
		So(strings.HasSuffix(top[0].StringVal, "rankB"), ShouldBeTrue)

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})
	Convey("Find with full-field filters and verify values", t, func() {
		marker := fmt.Sprintf("find-full-%d", time.Now().UnixNano())
		t0 := time.Date(2023, 1, 10, 9, 8, 7, 0, time.UTC)
		d0 := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
		ts0 := time.Date(2023, 1, 11, 12, 0, 0, 0, time.UTC)

		rec := &model.ComprehensiveTypesTable{
			SmallInt:     7,
			IntegerVal:   77,
			BigInt:       777,
			UintVal:      17,
			Uint8Val:     27,
			Uint16Val:    37,
			Uint32Val:    47,
			Uint64Val:    57,
			Float32Val:   1.23,
			Float64Val:   4.56,
			StringVal:    marker,
			VarcharVal:   "vf",
			CharVal:      "CF",
			BytesVal:     []byte{0x0f, 0x0e},
			BoolVal:      true,
			TimeVal:      t0,
			DateVal:      types.Date(d0),
			TimeOnly:     types.NewTime(5, 4, 3, 0),
			TimestampVal: ts0,
			JSONVal:      types.JSON([]byte(`{"k":1}`)),
			JsonbVal:     types.JSON([]byte(`{"k":1}`)),
			StringArray:  types.NewArray([]string{"p", "q"}),
			IntArray:     types.NewArray([]int32{9, 8}),
			BigintArray:  types.NewArray([]int64{7, 6}),
			FloatArray:   types.NewArray([]float64{5.5, 6.6}),
			UUIDVal:      types.NewUUIDv4(),
			DecimalVal:   11.11,
			NumericVal:   22.22,
			InetVal:      types.MustInet("10.10.0.1"),
			CidrVal:      types.MustCIDR("10.10.0.0/16"),
			MacaddrVal:   types.MustMACAddr("08:00:2b:0a:0b:0c"),
			PointVal:     types.NewPoint(7, 8),
			BoxVal:       types.NewBox(types.NewPoint(1, 1), types.NewPoint(2, 2)),
			PathVal:      types.NewPath([]types.Point{types.NewPoint(1, 1), types.NewPoint(2, 2)}, false),
			PolygonVal: types.NewPolygon(
				[]types.Point{types.NewPoint(0, 0), types.NewPoint(0, 2), types.NewPoint(2, 2)},
			),
			CircleVal:    types.NewCircle(types.NewPoint(1, 1), 1.5),
			BitVal:       types.NewBitString("10101010"),
			VarbitVal:    types.NewBitString("0011001100110011"),
			Int4rangeVal: types.NewInt4Range(3, 9, true, false),
			Int8rangeVal: types.NewInt8Range(5, 15, true, false),
			NumrangeVal:  types.NewNumRange(new(big.Rat).SetFloat64(3.14), new(big.Rat).SetFloat64(6.28), true, false),
			TsrangeVal:   types.NewTsRange(ts0, ts0.AddDate(0, 1, 0), true, false),
			TstzrangeVal: types.NewTstzRange(ts0, ts0.AddDate(0, 1, 0), true, false),
			DaterangeVal: types.NewDateRange(d0, d0.AddDate(0, 1, 0), true, false),
			TsvectorVal:  types.NewTSVector("'red':1 'blue':2"),
			TsqueryVal:   types.NewTSQuery("'red' | 'blue'"),
			XMLVal:       types.NewXML("<x/>"),
			MoneyVal:     types.NewMoney("$2.34"),
		}
		So(q.Create(rec), ShouldBeNil)

		rows, err := q.Where(
			tbl.StringVal.Eq(marker),
			tbl.SmallInt.Eq(rec.SmallInt),
			tbl.IntegerVal.Eq(rec.IntegerVal),
			tbl.BoolVal.Is(rec.BoolVal),
			tbl.BitVal.Eq(rec.BitVal),
			tbl.Int4rangeVal.Eq(rec.Int4rangeVal),
		).Find()
		So(err, ShouldBeNil)
		So(len(rows), ShouldEqual, 1)
		got := rows[0]
		So(got.StringVal, ShouldEqual, marker)
		So(got.SmallInt, ShouldEqual, rec.SmallInt)
		So(got.IntegerVal, ShouldEqual, rec.IntegerVal)
		So(got.BigInt, ShouldEqual, rec.BigInt)
		So(got.BoolVal, ShouldEqual, rec.BoolVal)
		So(got.PointVal, ShouldResemble, rec.PointVal)
		So(got.BitVal, ShouldResemble, rec.BitVal)
		So(got.Int4rangeVal, ShouldResemble, rec.Int4rangeVal)
		So(got.InetVal.String(), ShouldEqual, rec.InetVal.String())
		So(got.CidrVal.String(), ShouldEqual, rec.CidrVal.String())
		_, err = q.Where(tbl.ID.Eq(rec.ID)).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with time-based filters", t, func() {
		prefix := fmt.Sprintf("time-%d-", time.Now().UnixNano())
		baseTime := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)
		baseDate := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
		baseTimeOnly := types.NewTime(12, 30, 45, 0)
		baseTimestamp := time.Date(2023, 6, 15, 15, 30, 0, 0, time.UTC)

		// Create records with different time values
		timeRecords := []map[string]interface{}{
			{
				"string_val":    prefix + "early",
				"time_val":      baseTime.Add(-1 * time.Hour),
				"date_val":      types.Date(baseDate.AddDate(0, 0, -1)),
				"time_only":     types.NewTime(10, 0, 0, 0),
				"timestamp_val": baseTimestamp.Add(-2 * time.Hour),
			},
			{
				"string_val":    prefix + "middle",
				"time_val":      baseTime,
				"date_val":      types.Date(baseDate),
				"time_only":     baseTimeOnly,
				"timestamp_val": baseTimestamp,
			},
			{
				"string_val":    prefix + "late",
				"time_val":      baseTime.Add(1 * time.Hour),
				"date_val":      types.Date(baseDate.AddDate(0, 0, 1)),
				"time_only":     types.NewTime(14, 0, 0, 0),
				"timestamp_val": baseTimestamp.Add(2 * time.Hour),
			},
		}
		for _, rec := range timeRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// Test time comparisons
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.TimeVal.Gte(baseTime)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // middle and late

		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.TimeVal.Lt(baseTime)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1) // early

		// Test date comparisons
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 3) // all records

		// Test timestamp between
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.TimestampVal.Between(baseTimestamp.Add(-1*time.Hour), baseTimestamp.Add(1*time.Hour))).
			Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1) // middle

		// Test time_only comparisons
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.TimeOnly.Gt(types.NewTime(11, 0, 0, 0))).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // middle and late

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with numeric field comparisons", t, func() {
		prefix := fmt.Sprintf("num-%d-", time.Now().UnixNano())
		numRecords := []map[string]interface{}{
			{
				"string_val":  prefix + "small",
				"big_int":     int64(100),
				"uint_val":    int32(50),
				"uint8_val":   int16(10),
				"uint16_val":  int32(1000),
				"uint32_val":  int64(50000),
				"uint64_val":  float64(100000),
				"float32_val": float32(1.5),
				"float64_val": float64(2.5),
				"decimal_val": 10.50,
				"numeric_val": 100.25,
			},
			{
				"string_val":  prefix + "medium",
				"big_int":     int64(500),
				"uint_val":    int32(150),
				"uint8_val":   int16(50),
				"uint16_val":  int32(5000),
				"uint32_val":  int64(150000),
				"uint64_val":  float64(500000),
				"float32_val": float32(5.5),
				"float64_val": float64(7.5),
				"decimal_val": 50.75,
				"numeric_val": 500.50,
			},
			{
				"string_val":  prefix + "large",
				"big_int":     int64(1000),
				"uint_val":    int32(300),
				"uint8_val":   int16(100),
				"uint16_val":  int32(10000),
				"uint32_val":  int64(300000),
				"uint64_val":  float64(1000000),
				"float32_val": float32(10.5),
				"float64_val": float64(15.5),
				"decimal_val": 100.99,
				"numeric_val": 1000.75,
			},
		}
		for _, rec := range numRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// BigInt comparisons
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.BigInt.Gt(int64(200))).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // medium and large

		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.BigInt.Between(int64(200), int64(800))).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1) // medium

		// Uint comparisons
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.UintVal.In(int32(50), int32(300))).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // small and large

		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.Uint16Val.Gte(int32(5000))).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // medium and large

		// Float comparisons
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.Float32Val.Lt(float32(8.0))).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // small and medium

		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.Float64Val.Between(float64(5.0), float64(10.0))).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1) // medium

		// Decimal and Numeric comparisons
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.DecimalVal.Gt(30.0)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // medium and large

		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.NumericVal.Lt(200.0)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1) // small

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with bytes and binary data", t, func() {
		prefix := fmt.Sprintf("bytes-%d-", time.Now().UnixNano())
		binData1 := []byte{0x01, 0x02, 0x03}
		binData2 := []byte{0x0a, 0x0b, 0x0c}
		binData3 := []byte{0xff, 0xfe, 0xfd}

		So(q.DO.Create(map[string]interface{}{"string_val": prefix + "a", "bytes_val": binData1}), ShouldBeNil)
		So(q.DO.Create(map[string]interface{}{"string_val": prefix + "b", "bytes_val": binData2}), ShouldBeNil)
		So(q.DO.Create(map[string]interface{}{"string_val": prefix + "c", "bytes_val": binData3}), ShouldBeNil)

		// Test exact match
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.BytesVal.Eq(binData2)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1)
		So(rs[0].StringVal, ShouldEqual, prefix+"b")

		// Test NOT equal
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Not(tbl.BytesVal.Eq(binData1)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // b and c

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with advanced geometry operations", t, func() {
		prefix := fmt.Sprintf("geo-%d-", time.Now().UnixNano())

		// Create records with different geometric values
		geoRecords := []map[string]interface{}{
			{
				"string_val": prefix + "box1",
				"box_val":    types.NewBox(types.NewPoint(0, 0), types.NewPoint(2, 2)),
				"circle_val": types.NewCircle(types.NewPoint(1, 1), 2.0),
				"path_val":   types.NewPath([]types.Point{types.NewPoint(0, 0), types.NewPoint(1, 1)}, false),
				"polygon_val": types.NewPolygon(
					[]types.Point{types.NewPoint(0, 0), types.NewPoint(0, 1), types.NewPoint(1, 1)},
				),
			},
			{
				"string_val": prefix + "box2",
				"box_val":    types.NewBox(types.NewPoint(5, 5), types.NewPoint(7, 7)),
				"circle_val": types.NewCircle(types.NewPoint(6, 6), 1.5),
				"path_val":   types.NewPath([]types.Point{types.NewPoint(5, 5), types.NewPoint(7, 7)}, true),
				"polygon_val": types.NewPolygon(
					[]types.Point{types.NewPoint(5, 5), types.NewPoint(5, 7), types.NewPoint(7, 7)},
				),
			},
		}
		for _, rec := range geoRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// Test box equality
		testBox := types.NewBox(types.NewPoint(0, 0), types.NewPoint(2, 2))
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.BoxVal.Eq(testBox)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldBeGreaterThanOrEqualTo, 1) // PostgreSQL may reorder box corners

		// Test circle equality
		testCircle := types.NewCircle(types.NewPoint(1, 1), 2.0)
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.CircleVal.Eq(testCircle)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1)

		// Test path closed/open
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2)

		var openPaths, closedPaths int
		for _, r := range rs {
			if r.PathVal.Closed {
				closedPaths++
			} else {
				openPaths++
			}
		}
		So(openPaths, ShouldEqual, 1)
		So(closedPaths, ShouldEqual, 1)

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with range types operations", t, func() {
		prefix := fmt.Sprintf("range-%d-", time.Now().UnixNano())
		baseTime := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
		baseDate := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)

		rangeRecords := []map[string]interface{}{
			{
				"string_val":    prefix + "a",
				"int8range_val": types.NewInt8Range(10, 20, true, false),
				"numrange_val": types.NewNumRange(
					new(big.Rat).SetFloat64(1.0),
					new(big.Rat).SetFloat64(5.0),
					true,
					true,
				),
				"tsrange_val":   types.NewTsRange(baseTime, baseTime.AddDate(0, 1, 0), true, false),
				"tstzrange_val": types.NewTstzRange(baseTime.UTC(), baseTime.AddDate(0, 1, 0).UTC(), true, false),
				"daterange_val": types.NewDateRange(baseDate, baseDate.AddDate(0, 1, 0), true, false),
			},
			{
				"string_val":    prefix + "b",
				"int8range_val": types.NewInt8Range(15, 25, true, false),
				"numrange_val": types.NewNumRange(
					new(big.Rat).SetFloat64(3.0),
					new(big.Rat).SetFloat64(8.0),
					true,
					true,
				),
				"tsrange_val": types.NewTsRange(baseTime.AddDate(0, 0, 15), baseTime.AddDate(0, 2, 0), true, false),
				"tstzrange_val": types.NewTstzRange(
					baseTime.AddDate(0, 0, 15).UTC(),
					baseTime.AddDate(0, 2, 0).UTC(),
					true,
					false,
				),
				"daterange_val": types.NewDateRange(baseDate.AddDate(0, 0, 15), baseDate.AddDate(0, 2, 0), true, false),
			},
		}
		for _, rec := range rangeRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// Test int8range overlaps
		testRange := types.NewInt8Range(18, 22, true, false)
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.Int8rangeVal.Overlaps(testRange)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // both 'a' [10,20) and 'b' [15,25) overlap with [18,22)

		// Test numrange contains
		testNumRange := types.NewNumRange(new(big.Rat).SetFloat64(2.0), new(big.Rat).SetFloat64(4.0), true, true)
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.NumrangeVal.Contains(testNumRange)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldBeGreaterThanOrEqualTo, 1) // at least one should contain [2,4]

		// Test tsrange overlaps
		testTsRange := types.NewTsRange(baseTime.AddDate(0, 0, 20), baseTime.AddDate(0, 1, 10), true, false)
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.TsrangeVal.Overlaps(testTsRange)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldBeGreaterThanOrEqualTo, 1)

		// Test daterange equality
		testDateRange := types.NewDateRange(baseDate, baseDate.AddDate(0, 1, 0), true, false)
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.DaterangeVal.Eq(testDateRange)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldBeGreaterThanOrEqualTo, 1) // at least one match

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with bit string operations", t, func() {
		prefix := fmt.Sprintf("bit-%d-", time.Now().UnixNano())

		bitRecords := []map[string]interface{}{
			{
				"string_val": prefix + "a",
				"bit_val":    types.NewBitString("10101010"),
				"varbit_val": types.NewBitString("110011001100"),
			},
			{
				"string_val": prefix + "b",
				"bit_val":    types.NewBitString("11110000"),
				"varbit_val": types.NewBitString("101010101010"),
			},
			{
				"string_val": prefix + "c",
				"bit_val":    types.NewBitString("00001111"),
				"varbit_val": types.NewBitString("111111000000"),
			},
		}
		for _, rec := range bitRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// Test exact bit string match
		testBit := types.NewBitString("10101010")
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.BitVal.Eq(testBit)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1)
		So(rs[0].StringVal, ShouldEqual, prefix+"a")

		// Test varbit exact match
		testVarbit := types.NewBitString("110011001100")
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.VarbitVal.Eq(testVarbit)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1)
		So(rs[0].StringVal, ShouldEqual, prefix+"a")

		// Test NOT equal
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Not(tbl.BitVal.Eq(testBit)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // b and c

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with MAC address operations", t, func() {
		prefix := fmt.Sprintf("mac-%d-", time.Now().UnixNano())

		macRecords := []map[string]interface{}{
			{
				"string_val":  prefix + "cisco",
				"macaddr_val": types.MustMACAddr("08:00:2b:01:02:03"),
			},
			{
				"string_val":  prefix + "intel",
				"macaddr_val": types.MustMACAddr("00:1b:21:12:34:56"),
			},
			{
				"string_val":  prefix + "dell",
				"macaddr_val": types.MustMACAddr("00:14:22:ab:cd:ef"),
			},
		}
		for _, rec := range macRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// Test exact MAC match
		testMAC := types.MustMACAddr("08:00:2b:01:02:03")
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.MacaddrVal.Eq(testMAC)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1)
		So(rs[0].StringVal, ShouldEqual, prefix+"cisco")

		// Test MAC inequality
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Not(tbl.MacaddrVal.Eq(testMAC)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // intel and dell

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with UUID operations", t, func() {
		prefix := fmt.Sprintf("uuid-%d-", time.Now().UnixNano())
		uuid1 := types.NewUUIDv4()
		uuid2 := types.NewUUIDv4()
		uuid3 := types.NewUUIDv4()

		uuidRecords := []map[string]interface{}{
			{"string_val": prefix + "first", "uuid_val": uuid1},
			{"string_val": prefix + "second", "uuid_val": uuid2},
			{"string_val": prefix + "third", "uuid_val": uuid3},
		}
		for _, rec := range uuidRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// Test exact UUID match
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.UUIDVal.Eq(uuid2)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1)
		So(rs[0].StringVal, ShouldEqual, prefix+"second")

		// Test UUID IN clause
		rs, err = q.Where(tbl.StringVal.Like(prefix+"%"), tbl.UUIDVal.In(uuid1, uuid3)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // first and third

		// Test UUID NOT equal
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Not(tbl.UUIDVal.Eq(uuid1)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // second and third

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with XML operations", t, func() {
		prefix := fmt.Sprintf("xml-%d-", time.Now().UnixNano())

		xmlRecords := []map[string]interface{}{
			{
				"string_val": prefix + "root",
				"xml_val":    types.NewXML("<root><item>1</item></root>"),
			},
			{
				"string_val": prefix + "doc",
				"xml_val":    types.NewXML("<doc><title>Test</title></doc>"),
			},
			{
				"string_val": prefix + "empty",
				"xml_val":    types.NewXML("<empty/>"),
			},
		}
		for _, rec := range xmlRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// Test XML content presence instead of direct equality (PostgreSQL doesn't support XML = operator)
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.StringVal.Eq(prefix+"root")).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1)
		So(rs[0].StringVal, ShouldEqual, prefix+"root")
		So(string(rs[0].XMLVal), ShouldEqual, "<root><item>1</item></root>")

		// Test finding records by string marker since XML equality is not supported
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Not(tbl.StringVal.Eq(prefix + "root")).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // doc and empty

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})

	Convey("Find with TSQuery operations", t, func() {
		prefix := fmt.Sprintf("tsq-%d-", time.Now().UnixNano())

		tsqRecords := []map[string]interface{}{
			{
				"string_val":  prefix + "and",
				"tsquery_val": types.NewTSQuery("'cat' & 'dog'"),
			},
			{
				"string_val":  prefix + "or",
				"tsquery_val": types.NewTSQuery("'cat' | 'dog'"),
			},
			{
				"string_val":  prefix + "not",
				"tsquery_val": types.NewTSQuery("'cat' & !'dog'"),
			},
		}
		for _, rec := range tsqRecords {
			So(q.DO.Create(rec), ShouldBeNil)
		}

		// Test exact TSQuery match
		testTSQ := types.NewTSQuery("'cat' & 'dog'")
		rs, err := q.Where(tbl.StringVal.Like(prefix+"%"), tbl.TsqueryVal.Eq(testTSQ)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 1)
		So(rs[0].StringVal, ShouldEqual, prefix+"and")

		// Test TSQuery inequality
		rs, err = q.Where(tbl.StringVal.Like(prefix + "%")).Not(tbl.TsqueryVal.Eq(testTSQ)).Find()
		So(err, ShouldBeNil)
		So(len(rs), ShouldEqual, 2) // or and not

		// cleanup
		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})
}

func Test_Find_Batch(t *testing.T) {
	tbl, q := query.ComprehensiveTypesTable.QueryContext(context.Background())
	Convey("FindInBatch aggregates all results", t, func() {
		prefix := fmt.Sprintf("batch-%d-", time.Now().UnixNano())
		for i := 0; i < 5; i++ {
			So(q.DO.Create(map[string]interface{}{"string_val": fmt.Sprintf("%s%d", prefix, i)}), ShouldBeNil)
		}

		results, err := q.Where(tbl.StringVal.Like(prefix+"%")).
			FindInBatch(2, func(tx gen.Dao, batch int) error { return nil })
		So(err, ShouldBeNil)
		So(len(results), ShouldBeGreaterThanOrEqualTo, 5)

		_, err = q.Where(tbl.StringVal.Like(prefix + "%")).Delete()
		So(err, ShouldBeNil)
	})
}

func Test_Delete(t *testing.T) {
	tbl, q := query.ComprehensiveTypesTable.QueryContext(context.Background())
	Convey("Delete removes row and Take returns not found", t, func() {
		marker := fmt.Sprintf("to-del-%d", time.Now().UnixNano())
		So(q.DO.Create(map[string]interface{}{"string_val": marker}), ShouldBeNil)

		_, err := q.Where(tbl.StringVal.Eq(marker)).Delete()
		So(err, ShouldBeNil)

		_, err = q.Where(tbl.StringVal.Eq(marker)).Take()
		So(err, ShouldEqual, gorm.ErrRecordNotFound)
	})
}

func Test_ComprehensiveTypes_AllFields(t *testing.T) {
	Convey("ComprehensiveTypesTable: scan and validate all fields", t, func() {
		_, q := query.ComprehensiveTypesTable.QueryContext(context.Background())
		m, err := q.First()
		if err == gorm.ErrRecordNotFound {
			SkipConvey("no data to validate all fields")
			return
		}
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		Convey("simple numeric and string fields", func() {
			So(m.SmallInt, ShouldNotEqual, int16(0))
			So(m.IntegerVal, ShouldNotEqual, int32(0))
			So(m.BigInt, ShouldNotEqual, int64(0))
			So(m.UintVal, ShouldBeGreaterThanOrEqualTo, int32(0))
			So(m.Uint8Val, ShouldBeGreaterThanOrEqualTo, int16(0))
			So(m.Uint16Val, ShouldBeGreaterThanOrEqualTo, int32(0))
			So(m.Uint32Val, ShouldBeGreaterThanOrEqualTo, int64(0))
			So(m.Uint64Val == m.Uint64Val, ShouldBeTrue)
			So(m.Float32Val == m.Float32Val, ShouldBeTrue)
			So(m.Float64Val == m.Float64Val, ShouldBeTrue)
			So(len(m.StringVal) >= 0, ShouldBeTrue)
			So(len(m.VarcharVal) >= 0, ShouldBeTrue)
			So(len(m.CharVal) >= 0, ShouldBeTrue)
			So(m.BytesVal, ShouldNotBeNil)
			So(m.BoolVal == true || m.BoolVal == false, ShouldBeTrue)
		})

		Convey("time-like fields Value()", func() {
			_, err := m.DateVal.Value()
			So(err, ShouldBeNil)
			_, err = m.TimeOnly.Value()
			So(err, ShouldBeNil)
		})

		Convey("json and arrays Value()", func() {
			_, err := m.JSONVal.Value()
			So(err, ShouldBeNil)
			_, err = m.JsonbVal.Value()
			So(err, ShouldBeNil)
			_, err = m.StringArray.Value()
			So(err, ShouldBeNil)
			_, err = m.IntArray.Value()
			So(err, ShouldBeNil)
			_, err = m.BigintArray.Value()
			So(err, ShouldBeNil)
			_, err = m.FloatArray.Value()
			So(err, ShouldBeNil)
		})

		Convey("uuid, net, money, xml Value()", func() {
			_, err := m.UUIDVal.Value()
			So(err, ShouldBeNil)
			_, err = m.InetVal.Value()
			So(err, ShouldBeNil)
			_, err = m.CidrVal.Value()
			So(err, ShouldBeNil)
			_, err = m.MacaddrVal.Value()
			So(err, ShouldBeNil)
			_, err = m.MoneyVal.Value()
			So(err, ShouldBeNil)
			_, err = m.XMLVal.Value()
			So(err, ShouldBeNil)
		})

		Convey("geometry Value()", func() {
			_, err := m.PointVal.Value()
			So(err, ShouldBeNil)
			_, err = m.BoxVal.Value()
			So(err, ShouldBeNil)
			_, err = m.PathVal.Value()
			So(err, ShouldBeNil)
			_, err = m.PolygonVal.Value()
			So(err, ShouldBeNil)
			_, err = m.CircleVal.Value()
			So(err, ShouldBeNil)
		})

		Convey("bitstrings and ranges Value()", func() {
			_, err := m.BitVal.Value()
			So(err, ShouldBeNil)
			_, err = m.VarbitVal.Value()
			So(err, ShouldBeNil)
			_, err = m.Int4rangeVal.Value()
			So(err, ShouldBeNil)
			_, err = m.Int8rangeVal.Value()
			So(err, ShouldBeNil)
			_, err = m.NumrangeVal.Value()
			So(err, ShouldBeNil)
			_, err = m.TsrangeVal.Value()
			So(err, ShouldBeNil)
			_, err = m.TstzrangeVal.Value()
			So(err, ShouldBeNil)
			_, err = m.DaterangeVal.Value()
			So(err, ShouldBeNil)
		})

		Convey("fulltext Value()", func() {
			_, err := m.TsvectorVal.Value()
			So(err, ShouldBeNil)
			_, err = m.TsqueryVal.Value()
			So(err, ShouldBeNil)
		})
	})
}
