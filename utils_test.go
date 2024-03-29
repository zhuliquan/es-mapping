package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckTypeSupportLucene(t *testing.T) {
	type args struct {
		typ FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_true",
			args: args{typ: KEYWORD_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_false",
			args: args{typ: SHAPE_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckTypeSupportLucene(tt.args.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractFieldAliasMap(t *testing.T) {
	type args struct {
		m *Mapping
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "test_error_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: ALIAS_FIELD_TYPE,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_02",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"first": {
										Type: ALIAS_FIELD_TYPE,
									},
								},
							},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_03",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"first": {
										Type: ALIAS_FIELD_TYPE,
										Path: "name.first",
									},
								},
							},
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_error_04",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"alias_type": {
							Type: ALIAS_FIELD_TYPE,
							Path: "other_field",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test_ok_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name_2": {
							Type: ALIAS_FIELD_TYPE,
							Path: "name_second",
						},
						"name_second": {
							Type: TEXT_FIELD_TYPE,
						},
					},
				},
			},
			want:    map[string]string{"name_2": "name_second"},
			wantErr: false,
		},
		{
			name: "test_ok_02",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"first": {
										Type: ALIAS_FIELD_TYPE,
										Path: "name_first",
									},
								},
							},
						},
						"name_first": {
							Type: TEXT_FIELD_TYPE,
						},
					},
				},
			},
			want:    map[string]string{"name.first": "name_first"},
			wantErr: false,
		},
		{
			name: "test_ok_03",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"name": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"first": {
										Type: TEXT_FIELD_TYPE,
									},
								},
							},
						},
						"name_first": {
							Type: ALIAS_FIELD_TYPE,
							Path: "name.first",
						},
					},
				},
			},
			want:    map[string]string{"name_first": "name.first"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pm = &PropertyMapping{
				fieldMapping:  tt.args.m,
				propertyCache: map[string]*Property{},
				fieldAliasMap: map[string]string{},
			}
			got, err := extractFieldAliasMap(pm)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFetchProperty(t *testing.T) {
	type args struct {
		m      *Mapping
		target string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*Property
		wantErr bool
	}{
		{
			name: "test_error_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"y": {Type: TEXT_FIELD_TYPE},
					},
				},
				target: "x",
			},
			want: map[string]*Property{},
		},
		{
			name: "test_error_02",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: TEXT_FIELD_TYPE,
						},
					},
				},
				target: "x.y",
			},
			want: map[string]*Property{},
		},
		{
			name: "test_error_03",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y": {
										Type: TEXT_FIELD_TYPE,
									},
								},
							},
						},
					},
				},
				target: "x",
			},
			want: map[string]*Property{},
		},
		{
			name: "test_error_04",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y": {
										Type: OBJECT_FIELD_TYPE,
										Mapping: Mapping{
											Properties: map[string]*Property{
												"z": {
													Type: TEXT_FIELD_TYPE,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				target: "x.y",
			},
			want: map[string]*Property{},
		},
		{
			name: "test_err_05",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y": {
										Type: OBJECT_FIELD_TYPE,
										Mapping: Mapping{
											Properties: map[string]*Property{
												"z.a": {
													Type: TEXT_FIELD_TYPE,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				target: "x.y.z",
			},
			want: map[string]*Property{},
		},
		{
			name: "test_ok_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: FLATTENED_FIELD_TYPE,
						},
					},
				},
				target: "x",
			},
			want: map[string]*Property{},
		},
		{
			name: "test_ok_01_01",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: FLATTENED_FIELD_TYPE,
						},
					},
				},
				target: "x.y",
			},
			want: map[string]*Property{"x.y": {Type: KEYWORD_FIELD_TYPE}},
		},
		{
			name: "test_ok_02",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y": {
										Type: TEXT_FIELD_TYPE,
									},
								},
							},
						},
					},
				},
				target: "x.y",
			},
			want: map[string]*Property{
				"x.y": {
					Type: TEXT_FIELD_TYPE,
				},
			},
		},
		{
			name: "test_ok_03",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: TEXT_FIELD_TYPE,
							Fields: map[string]*Property{
								"raw": {
									Type: KEYWORD_FIELD_TYPE,
								},
							},
						},
					},
				},
				target: "x.raw",
			},
			want: map[string]*Property{
				"x.raw": {
					Type: KEYWORD_FIELD_TYPE,
				},
			},
		},
		{
			name: "test_ok_04",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: TEXT_FIELD_TYPE,
							Fields: map[string]*Property{
								"raw1": {
									Type: KEYWORD_FIELD_TYPE,
								},
								"raw2": {
									Type: LONG_FIELD_TYPE,
								},
							},
						},
					},
				},
				target: "x.raw\\*",
			},
			want: map[string]*Property{
				"x.raw1": {
					Type: KEYWORD_FIELD_TYPE,
				},
				"x.raw2": {
					Type: LONG_FIELD_TYPE,
				},
			},
		},
		{
			name: "test_err_05",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: FLATTENED_FIELD_TYPE,
						},
					},
				},
				target: "x.*",
			},
			want: map[string]*Property{},
		},
		{
			name: "test_err_06",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y.z": {Type: KEYWORD_FIELD_TYPE},
									"y":   {Type: FLATTENED_FIELD_TYPE},
								},
							},
						},
					},
				},
				target: "x.y.z",
			},
			want: map[string]*Property{
				"x.y.z": {Type: KEYWORD_FIELD_TYPE},
			},
		},
		{
			name: "test_err_07",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y.z": {Type: DATE_FIELD_TYPE},
									"y":   {Type: FLATTENED_FIELD_TYPE},
								},
							},
						},
					},
				},
				target: "x.y.z",
			},
			want:    map[string]*Property{},
			wantErr: true,
		},
		{
			name: "test_err_08",
			args: args{
				m: &Mapping{
					Properties: map[string]*Property{
						"x": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Properties: map[string]*Property{
									"y":   {Type: FLATTENED_FIELD_TYPE},
									"y.z": {Type: DATE_FIELD_TYPE},
								},
							},
						},
					},
				},
				target: "x.y.z",
			},
			want:    map[string]*Property{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m = &PropertyMapping{
				fieldMapping: tt.args.m,
			}
			got, err := getProperty(m, tt.args.target)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCheckIntType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_check_int_01",
			args: args{t: INTEGER_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_int_02",
			args: args{t: INTEGER_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_byte",
			args: args{t: BYTE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_short",
			args: args{t: SHORT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_long_01",
			args: args{t: LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_long_02",
			args: args{t: LONG_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_other",
			args: args{t: DOUBLE_RANGE_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckIntType(tt.args.t); got != tt.want {
				t.Errorf("CheckIntType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckUIntType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_uint64",
			args: args{t: UNSIGNED_LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_other",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckUIntType(tt.args.t); got != tt.want {
				t.Errorf("CheckUIntType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckFloatType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_float16",
			args: args{t: HALF_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float32_01",
			args: args{t: FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float32_02",
			args: args{t: FLOAT_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float64_01",
			args: args{t: DOUBLE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float64_02",
			args: args{t: DOUBLE_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float128",
			args: args{t: SCALED_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_other",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckFloatType(tt.args.t); got != tt.want {
				t.Errorf("CheckFloatType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDateType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_check_date_01",
			args: args{t: DATE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_date_02",
			args: args{t: DATE_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_date_03",
			args: args{t: DATE_NANOS_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_other",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDateType(tt.args.t); got != tt.want {
				t.Errorf("CheckDateType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckNumberType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_check_int_01",
			args: args{t: INTEGER_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_int_02",
			args: args{t: INTEGER_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_byte",
			args: args{t: BYTE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_short",
			args: args{t: SHORT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_long_01",
			args: args{t: LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_long_02",
			args: args{t: LONG_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_uint64",
			args: args{t: UNSIGNED_LONG_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float16",
			args: args{t: HALF_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float32_01",
			args: args{t: FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float32_02",
			args: args{t: FLOAT_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float64_01",
			args: args{t: DOUBLE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float64_02",
			args: args{t: DOUBLE_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_float128",
			args: args{t: SCALED_FLOAT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_check_other",
			args: args{t: FLATTENED_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckNumberType(tt.args.t); got != tt.want {
				t.Errorf("CheckNumberType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckIPType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_ip_01",
			args: args{t: IP_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_ip_02",
			args: args{t: IP_RANGE_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_other",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckIPType(tt.args.t); got != tt.want {
				t.Errorf("CheckIPType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckVersionType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_version",
			args: args{t: VERSION_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_other_version",
			args: args{t: UNKNOWN_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckVersionType(tt.args.t); got != tt.want {
				t.Errorf("CheckVersionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckStringType(t *testing.T) {
	type args struct {
		t FieldType
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_keyword",
			args: args{t: KEYWORD_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_text",
			args: args{t: TEXT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_wildcard",
			args: args{t: WILDCARD_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_constant_keyword",
			args: args{t: CONSTANT_KEYWORD_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_match_only_text",
			args: args{t: MATCH_ONLY_TEXT_FIELD_TYPE},
			want: true,
		},
		{
			name: "test_int",
			args: args{t: INTEGER_FIELD_TYPE},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckStringType(tt.args.t); got != tt.want {
				t.Errorf("CheckVersionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFillDefaultParameter(t *testing.T) {
	type args struct {
		pm *PropertyMapping
	}
	tests := []struct {
		name string
		args args
		want *PropertyMapping
	}{
		{
			name: "fill_date_scaling_type",
			args: args{pm: &PropertyMapping{
				fieldMapping: &Mapping{
					Properties: map[string]*Property{
						"date1": {
							Type: DATE_FIELD_TYPE,
						},
						"date2": {
							Type: DATE_RANGE_FIELD_TYPE,
						},
						"date3": {
							Type: DATE_NANOS_FIELD_TYPE,
						},
						"scaling": {
							Type: SCALED_FLOAT_FIELD_TYPE,
						},
						"nest1": {
							Mapping: Mapping{
								Properties: map[string]*Property{
									"date1": {
										Type: DATE_FIELD_TYPE,
									},
									"date2": {
										Type: DATE_RANGE_FIELD_TYPE,
									},
									"date3": {
										Type: DATE_NANOS_FIELD_TYPE,
									},
									"scaling": {
										Type: SCALED_FLOAT_FIELD_TYPE,
									},
								},
							},
						},
					},
				},
			}},
			want: &PropertyMapping{
				fieldMapping: &Mapping{
					Dynamic: BoolDynamic(true),
					Properties: map[string]*Property{
						"date1": {
							Type:   DATE_FIELD_TYPE,
							Format: "strict_date_optional_time||epoch_millis",
						},
						"date2": {
							Type:   DATE_RANGE_FIELD_TYPE,
							Format: "strict_date_optional_time||epoch_millis",
						},
						"date3": {
							Type:   DATE_NANOS_FIELD_TYPE,
							Format: "strict_date_optional_time_nanos||epoch_millis",
						},
						"scaling": {
							Type:          SCALED_FLOAT_FIELD_TYPE,
							ScalingFactor: 1.0,
						},
						"nest1": {
							Type: OBJECT_FIELD_TYPE,
							Mapping: Mapping{
								Dynamic: BoolDynamic(true),
								Properties: map[string]*Property{
									"date1": {
										Type:   DATE_FIELD_TYPE,
										Format: "strict_date_optional_time||epoch_millis",
									},
									"date2": {
										Type:   DATE_RANGE_FIELD_TYPE,
										Format: "strict_date_optional_time||epoch_millis",
									},
									"date3": {
										Type:   DATE_NANOS_FIELD_TYPE,
										Format: "strict_date_optional_time_nanos||epoch_millis",
									},
									"scaling": {
										Type:          SCALED_FLOAT_FIELD_TYPE,
										ScalingFactor: 1.0,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fillDefaultParameter(tt.args.pm)
			assert.Equal(t, tt.want, tt.args.pm)
		})
	}
}
