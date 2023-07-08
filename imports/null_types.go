package imports

import (
	"database/sql"
	"time"
)

func NullBoolPtr(Bool *bool) sql.NullBool {
	if Bool != nil {
		return sql.NullBool{Bool: *Bool, Valid: true}
	} else {
		return sql.NullBool{}
	}
}

func NullFloat64Ptr(Float64 *float64) sql.NullFloat64 {
	if Float64 != nil {
		return sql.NullFloat64{Float64: *Float64, Valid: true}
	} else {
		return sql.NullFloat64{}
	}
}

func NullInt32Ptr(Int32 *int32) sql.NullInt32 {
	if Int32 != nil {
		return sql.NullInt32{Int32: *Int32, Valid: true}
	} else {
		return sql.NullInt32{}
	}
}

func NullInt16Ptr(Int16 *int16) sql.NullInt16 {
	if Int16 != nil {
		return sql.NullInt16{Int16: *Int16, Valid: true}
	} else {
		return sql.NullInt16{}
	}
}

func NullInt64Ptr(Int64 *int64) sql.NullInt64 {
	if Int64 != nil {
		return sql.NullInt64{*Int64, true}
	} else {
		return sql.NullInt64{0, false}
	}
}

func NullStringPtr(String *string) sql.NullString {
	if String != nil {
		return sql.NullString{String: *String, Valid: true}
	} else {
		return sql.NullString{}
	}
}

func NullTimePtr(Time *time.Time) sql.NullTime {
	if Time != nil {
		return sql.NullTime{Time: *Time, Valid: true}
	} else {
		return sql.NullTime{}
	}
}

func NullBoolToPtr(Bool sql.NullBool) *bool {
	if Bool.Valid {
		return &Bool.Bool
	} else {
		return nil
	}
}

func NullFloat64ToPtr(Float64 sql.NullFloat64) *float64 {
	if Float64.Valid {
		return &Float64.Float64
	} else {
		return nil
	}
}

func NullInt32ToPtr(Int32 sql.NullInt32) *int32 {
	if Int32.Valid {
		return &Int32.Int32
	} else {
		return nil
	}
}

func NullInt16ToPtr(Int16 sql.NullInt16) *int16 {
	if Int16.Valid {
		return &Int16.Int16
	} else {
		return nil
	}
}

func NullInt64ToPtr(Int64 sql.NullInt64) *int64 {
	if Int64.Valid {
		return &Int64.Int64
	} else {
		return nil
	}
}

func NullStringToPtr(String sql.NullString) *string {
	if String.Valid {
		return &String.String
	} else {
		return nil
	}
}

func NullTimeToPtr(Time sql.NullTime) *time.Time {
	if Time.Valid {
		return &Time.Time
	} else {
		return nil
	}
}
