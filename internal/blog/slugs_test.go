package blog

import "testing"

func TestToSlug(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "C# maps to c-sharp",
			tag:  "C#",
			want: "c-sharp",
		},
		{
			name: "c# lowercase maps to c-sharp",
			tag:  "c#",
			want: "c-sharp",
		},
		{
			name: "C++ maps to cpp",
			tag:  "C++",
			want: "cpp",
		},
		{
			name: "c++ lowercase maps to cpp",
			tag:  "c++",
			want: "cpp",
		},
		{
			name: "normal tag is slugified",
			tag:  "Software Engineering",
			want: "software-engineering",
		},
		{
			name: "single word tag",
			tag:  "Go",
			want: "go",
		},
		{
			name: "tag with special characters",
			tag:  "ASP.NET",
			want: "asp-net",
		},
		{
			name: "empty string",
			tag:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSlug(tt.tag)
			if got != tt.want {
				t.Errorf("ToSlug(%q) = %q, want %q", tt.tag, got, tt.want)
			}
		})
	}
}
