//+build integration

package integrations_test

import (
	. "friendly-potato/integrations"
	"os"
	"reflect"
	"testing"
)

const(
	realDomain1="1-qr.me"
	realDomain2="gitpages.dev"
	fakeDomain="abcxpto2020ultra.com"
	)

func TestInitAPI(t *testing.T) {
	t.Logf("lookup CF_TOKEN env")
	cfToken, ok := os.LookupEnv("CF_TOKEN")
	if !ok {
		t.Skipf("could not find cloudflare token as env SKIPPING...")
	}
	type args struct {
		apiToken string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test init cloudflare api",
			args:    args{cfToken},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitCloudFlareAPI(tt.args.apiToken); (err != nil) != tt.wantErr {
				t.Errorf("InitCloudFlareAPI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateZone(t *testing.T) {
	t.Logf("lookup CF_TOKEN env")
	cfToken, ok := os.LookupEnv("CF_TOKEN")
	if !ok {
		t.Skipf("could not find cloudflare token as env SKIPPING...")
	} else {
		err := InitCloudFlareAPI(cfToken)
		if err != nil {
			t.Logf("failed to initialize cloudflare api: %v", err)
		}
	}
	type args struct {
		zone Zone
	}
	tests := []struct {
		name            string
		args            args
		wantCreatedZone Zone
		wantErr         bool
	}{
		{
			name: "create zone "+realDomain1,
			args: args{
				zone: Zone{
					Resource: DomainResource{Name: realDomain1,},
				},
			},
			wantCreatedZone: Zone{
				Resource: DomainResource{Name: realDomain1,},
			},
			wantErr: false,
		},
		{
			name: "create zone "+realDomain2,
			args: args{
				zone: Zone{
					Resource: DomainResource{Name: realDomain2,},
				},
			},
			wantCreatedZone: Zone{
				Resource: DomainResource{Name: realDomain2,},
			},
			wantErr: false,
		},
		{
			name: "create duplicate zone " + realDomain2,
			args: args{
				zone: Zone{
					Resource: DomainResource{Name: realDomain2,},
				},
			},
			wantCreatedZone: Zone{
				Resource: DomainResource{Name: ""},
			},
			wantErr: true,
		},
		{
			name: "create fake domain zone " + fakeDomain,
			args: args{
				zone: Zone{
					Resource: DomainResource{Name: fakeDomain,},
				},
			},
			wantCreatedZone: Zone{
				Resource: DomainResource{Name: ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCreatedZone, err := CreateZone(tt.args.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCreatedZone.Resource.Name, tt.wantCreatedZone.Resource.Name) {
				t.Errorf("CreateZone() gotCreatedZone = %v, want %v", gotCreatedZone.Resource, tt.wantCreatedZone.Resource.Name)
			}
		})
	}
}

func TestListZones(t *testing.T) {
	t.Logf("lookup CF_TOKEN env")
	cfToken, ok := os.LookupEnv("CF_TOKEN")
	if !ok {
		t.Skipf("could not find cloudflare token as env SKIPPING...")
	} else {
		err := InitCloudFlareAPI(cfToken)
		if err != nil {
			t.Logf("failed to initialize cloudflare api: %v", err)
		}
	}
	tests := []struct {
		name      string
		wantZones Zones
		wantErr   bool
	}{
		{
			name: "list zone find zone "+ realDomain1 +","+ realDomain2,
			wantZones: Zones{
				Zone{
					Resource: DomainResource{Name: realDomain1},
				},
				{
					Resource: DomainResource{Name: realDomain2},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotZones, err := ListZones()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListZones() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotZones.Names(), tt.wantZones.Names()) {
				t.Errorf("ListZones() gotZones = %v, want %v", gotZones.Names(), tt.wantZones.Names())
			}
		})
	}
}

func TestDeleteZone(t *testing.T) {
	t.Logf("lookup CF_TOKEN env")
	cfToken, ok := os.LookupEnv("CF_TOKEN")
	if !ok {
		t.Skipf("could not find cloudflare token as env SKIPPING...")
	} else {
		err := InitCloudFlareAPI(cfToken)
		if err != nil {
			t.Logf("failed to initialize cloudflare api: %v", err)
		}
	}
	type args struct {
		zone Zone
	}
	tests := []struct {
		name            string
		args            args
		wantDeletedZone Zone
		wantErr         bool
	}{
		{
			name: "delete zone " + realDomain2,
			args: args{
				zone: Zone{
					Resource: DomainResource{
						Name: realDomain2,
					},
				},
			},
			wantDeletedZone: Zone{
				Resource: DomainResource{Name: realDomain2,},
			},
			wantErr: false,
		},
		{
			name: "delete zone "+realDomain1,
			args: args{
				zone: Zone{
					Resource: DomainResource{Name: realDomain1,},
				},
			},
			wantDeletedZone: Zone{
				Resource: DomainResource{Name: realDomain1,},
			},
			wantErr: false,
		},
		{
			name: "delete nonexistent zone" + fakeDomain,
			args: args{
				zone: Zone{
					Resource: DomainResource{Name: fakeDomain,},
				},
			},
			wantDeletedZone: Zone{
				Resource: DomainResource{Name: "",},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeletedZone, err := DeleteZone(tt.args.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDeletedZone.Resource.Name, tt.wantDeletedZone.Resource.Name) {
				t.Errorf("DeleteZone() gotDeletedZone = %v, want %v", gotDeletedZone.Resource.Name, tt.wantDeletedZone.Resource.Name)
			}
		})
	}
}
