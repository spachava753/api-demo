package domain

import (
	"reflect"
	"testing"
)

func stringPtr(s string) *string {
	return &s
}

func Test_matchesUser(t *testing.T) {

	type args struct {
		user       User
		properties UserProperties
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no properties",
			args: args{
				user: User{
					Name: "Shashank Pachava",
					Role: "admin",
				},
				properties: UserProperties{},
			},
			want: true,
		},
		{
			name: "name",
			args: args{
				user: User{
					Name: "Shashank Pachava",
					Role: "admin",
				},
				properties: UserProperties{
					Name: stringPtr("Shashank Pachava"),
				},
			},
			want: true,
		},
		{
			name: "role",
			args: args{
				user: User{
					Name: "Shashank Pachava",
					Role: "admin",
				},
				properties: UserProperties{
					Role: stringPtr("admin"),
				},
			},
			want: true,
		},
		{
			name: "name and role",
			args: args{
				user: User{
					Name: "Shashank Pachava",
					Role: "admin",
				},
				properties: UserProperties{
					Name: stringPtr("Shashank Pachava"),
					Role: stringPtr("admin"),
				},
			},
			want: true,
		},
		{
			name: "don't match name",
			args: args{
				user: User{
					Name: "Shashank Pachava",
					Role: "admin",
				},
				properties: UserProperties{
					Name: stringPtr("Shashank"),
				},
			},
			want: false,
		},
		{
			name: "don't match role",
			args: args{
				user: User{
					Name: "Shashank Pachava",
					Role: "admin",
				},
				properties: UserProperties{
					Role: stringPtr("user"),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesUser(tt.args.user, tt.args.properties); got != tt.want {
				t.Errorf("matchesUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterUsers(t *testing.T) {
	type args struct {
		users      []User
		properties UserProperties
	}
	tests := []struct {
		name string
		args args
		want []User
	}{
		{
			name: "no properties",
			args: args{
				users: []User{
					{
						Name: "Shashank Pachava",
						Role: "admin",
					},
				},
				properties: UserProperties{},
			},
			want: []User{
				{
					Name: "Shashank Pachava",
					Role: "admin",
				},
			},
		},
		{
			name: "role properties",
			args: args{
				users: []User{
					{
						Name: "Shashank Pachava",
						Role: "admin",
					},
					{
						Name: "Sasi",
						Role: "user",
					},
					{
						Name: "Sridhar",
						Role: "user",
					},
				},
				properties: UserProperties{
					Role: stringPtr("admin"),
				},
			},
			want: []User{
				{
					Name: "Shashank Pachava",
					Role: "admin",
				},
			},
		},
		{
			name: "role properties",
			args: args{
				users: []User{
					{
						Name: "Sasi",
						Role: "user",
					},
					{
						Name: "Shashank Pachava",
						Role: "admin",
					},
					{
						Name: "Sridhar",
						Role: "user",
					},
				},
				properties: UserProperties{
					Role: stringPtr("admin"),
				},
			},
			want: []User{
				{
					Name: "Shashank Pachava",
					Role: "admin",
				},
			},
		},
		{
			name: "role properties",
			args: args{
				users: []User{
					{
						Name: "Sasi",
						Role: "user",
					},
					{
						Name: "Sridhar",
						Role: "user",
					},
					{
						Name: "Shashank Pachava",
						Role: "admin",
					},
				},
				properties: UserProperties{
					Role: stringPtr("admin"),
				},
			},
			want: []User{
				{
					Name: "Shashank Pachava",
					Role: "admin",
				},
			},
		},
		{
			name: "role properties",
			args: args{
				users: []User{
					{
						Name: "Sasi",
						Role: "user",
					},
					{
						Name: "Sridhar",
						Role: "user",
					},
					{
						Name: "Shashank Pachava",
						Role: "admin",
					},
				},
				properties: UserProperties{
					Role: stringPtr("user"),
				},
			},
			want: []User{
				{
					Name: "Sasi",
					Role: "user",
				},
				{
					Name: "Sridhar",
					Role: "user",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterUsers(tt.args.users, tt.args.properties); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
