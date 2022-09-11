package array

import "testing"

type Sample struct {
    A string `array:"[0]"`
    C int `array:"[1]"`
    B bool `array:"[2]"`
}

type FailedSample struct {
    A string `array:"[0]"`
    B Sample `array:"[2]"`
}

type OmitSample struct {
    A string `array:"[0]"`
    B bool `array:"[100] ,omitempty"`
}

func TestParse(t *testing.T) {
    type args struct {
        d []string
        v *Sample
    }
    tests := []struct {
        name    string
        args    args
        want    *Sample
        wantErr bool
    }{
        {
            "",
            args{
                []string{"first", "2", "false"},
                &Sample{},
            },
            &Sample{"first",2,false},
            false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Parse(tt.args.d, tt.args.v)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
            } else if *tt.args.v != *tt.want {
                t.Errorf("Parse() get = %v, want %v", tt.args.v, tt.want)
            } else{
                t.Log(tt.args.v)
            }
        })
    }
}

func TestParseOmit(t *testing.T) {
    type args struct {
        d []string
        v *OmitSample
    }
    tests := []struct {
        name    string
        args    args
        want    *OmitSample
        wantErr bool
    }{
        {
            "",
            args{
                []string{"first", "second", "third"},
                &OmitSample{},
            },
            &OmitSample{"first", false},
            false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Parse(tt.args.d, tt.args.v)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
            } else if *tt.args.v != *tt.want {
                t.Errorf("Parse() get = %v, want %v", tt.args.v, tt.want)
            }
        })
    }
}


func TestParseFailed(t *testing.T) {
    type args struct {
        d []string
        v *FailedSample
    }
    tests := []struct {
        name    string
        args    args
        want    *FailedSample
        wantErr bool
    }{
        {
            "",
            args{
                []string{"first", "second", "third"},
                &FailedSample{},
            },
            &FailedSample{"first", Sample{}},
            true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Parse(tt.args.d, tt.args.v)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
            } else if *tt.args.v != *tt.want {
                t.Errorf("Parse() get = %v, want %v", tt.args.v, tt.want)
            }
        })
    }
}
