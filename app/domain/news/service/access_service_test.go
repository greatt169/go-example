package service

import (
	helpers "github.com/AeroAgency/golang-helpers-lib"
	"github.com/AeroAgency/golang-helpers-lib/dto"
	"news-ms/domain/news/entity"
	"testing"
	"time"
)

func TestNewNewsAccess(t *testing.T) {
	NewNewsAccess()
}

func TestNewsAccess_CheckCanShowDraftsNews(t *testing.T) {
	tests := []struct {
		name          string
		access        *helpers.Access
		privilegesDto dto.Privileges
		want          bool
	}{
		{
			name:   "Check privileges [show_deactivated] for user without scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_active"},
				},
			},
			want: false,
		},
		{
			name:   "Check privileges [show_deactivated] for user with scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_deactivated"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if got := a.CheckCanShowDraftsNews(tt.privilegesDto); got != tt.want {
				t.Errorf("CheckCanShowDraftsNews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsAccess_CheckCanShowDelayedNews(t *testing.T) {
	tests := []struct {
		name          string
		access        *helpers.Access
		privilegesDto dto.Privileges
		want          bool
	}{
		{
			name:   "Check privileges [show_delayed] for user without scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_active"},
				},
			},
			want: false,
		},
		{
			name:   "Check privileges [show_delayed] for user with scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_delayed"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if got := a.CheckCanShowDelayedNews(tt.privilegesDto); got != tt.want {
				t.Errorf("CheckCanShowDelayedNews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsAccess_CheckCanShowDraftNewsOrFail(t *testing.T) {
	type args struct {
		privilegesDto dto.Privileges
		newsEntity    *entity.News
	}
	tests := []struct {
		name    string
		access  *helpers.Access
		args    args
		wantErr bool
	}{
		{
			name:   "Check privileges [show_deactivated] for user without scope and entity [Active=false]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges [show_deactivated] for user without scope and entity [Active=true]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges [show_deactivated] for user with scope and entity [Active=false]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges [show_deactivated] for user with scope and entity [Active=true]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if err := a.CheckCanShowDraftNewsOrFail(tt.args.privilegesDto, tt.args.newsEntity); (err != nil) != tt.wantErr {
				t.Errorf("CheckCanShowDraftNewsOrFail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsAccess_CheckCanShowDelayedNewsOrFail(t *testing.T) {
	type args struct {
		privilegesDto dto.Privileges
		newsEntity    *entity.News
	}
	tests := []struct {
		name    string
		access  *helpers.Access
		args    args
		wantErr bool
	}{
		{
			name:   "Check privileges [show_delayed] for user without scope and entity [ActiveFrom > now()]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: &entity.News{
					ActiveFrom: time.Now().Unix() + 1,
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges [show_delayed] for user without scope and entity [ActiveFrom < now()]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: &entity.News{
					ActiveFrom: time.Now().Unix() - 1,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges [show_delayed] for user with scope and entity [ActiveFrom > now()]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_delayed"},
					},
				},
				newsEntity: &entity.News{
					ActiveFrom: time.Now().Unix() + 1,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges [show_delayed] for user with scope and entity [ActiveFrom < now()]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_delayed"},
					},
				},
				newsEntity: &entity.News{
					ActiveFrom: time.Now().Unix() - 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if err := a.CheckCanShowDelayedNewsOrFail(tt.args.privilegesDto, tt.args.newsEntity); (err != nil) != tt.wantErr {
				t.Errorf("CheckCanShowDelayedNewsOrFail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsAccess_CheckCanCreateNewsOrFail(t *testing.T) {
	tests := []struct {
		name          string
		access        *helpers.Access
		privilegesDto dto.Privileges
		wantErr       bool
	}{
		{
			name:   "Check privileges [create] for user without scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_active"},
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges [create] for user with scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"create"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if err := a.CheckCanCreateNewsOrFail(tt.privilegesDto); (err != nil) != tt.wantErr {
				t.Errorf("CheckCanCreateNewsOrFail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsAccess_checkCanUpdateActiveNews(t *testing.T) {
	tests := []struct {
		name          string
		access        *helpers.Access
		privilegesDto dto.Privileges
		want          bool
	}{
		{
			name:   "Check privileges [update_active] for user without scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_active"},
				},
			},
			want: false,
		},
		{
			name:   "Check privileges [update_active] for user with scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"update_active"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if got := a.checkCanUpdateActiveNews(tt.privilegesDto); got != tt.want {
				t.Errorf("checkCanUpdateActiveNews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsAccess_checkCanUpdateDeactivatedNews(t *testing.T) {
	tests := []struct {
		name          string
		access        *helpers.Access
		privilegesDto dto.Privileges
		want          bool
	}{
		{
			name:   "Check privileges [update_deactivated] for user without scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_active"},
				},
			},
			want: false,
		},
		{
			name:   "Check privileges [update_deactivated] for user with scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"update_deactivated"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if got := a.checkCanUpdateDeactivatedNews(tt.privilegesDto); got != tt.want {
				t.Errorf("checkCanUpdateDeactivatedNews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsAccess_CheckCanUpdateNewsOrFail(t *testing.T) {
	type args struct {
		privilegesDto dto.Privileges
		newsEntity    *entity.News
	}
	tests := []struct {
		name    string
		access  *helpers.Access
		args    args
		wantErr bool
	}{
		{
			name:   "Check privileges to update News [Active=true] for User with scopes [update_active, update_deactivated]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active", "update_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges to update News [Active=false] for User with scopes [update_active, update_deactivated]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active", "update_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges to update News [Active=true] for User with scopes [update_deactivated]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges to update News [Active=false] for User with scopes [update_deactivated]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges to update News [Active=true] for User with scopes [update_active]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges to update News [Active=false] for User with scopes [update_active]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"update_active"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges to update News [Active=true] for User without update scopes",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges to update News [Active=false] for User without update scopes",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if err := a.CheckCanUpdateNewsOrFail(tt.args.privilegesDto, tt.args.newsEntity); (err != nil) != tt.wantErr {
				t.Errorf("CheckCanUpdateNewsOrFail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsAccess_checkCanDeleteActiveNews(t *testing.T) {
	tests := []struct {
		name          string
		access        *helpers.Access
		privilegesDto dto.Privileges
		want          bool
	}{
		{
			name:   "Check privileges [delete_active] for user without scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_active"},
				},
			},
			want: false,
		},
		{
			name:   "Check privileges [delete_active] for user with scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"delete_active"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if got := a.checkCanDeleteActiveNews(tt.privilegesDto); got != tt.want {
				t.Errorf("checkCanDeleteActiveNews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsAccess_checkCanDeleteDeactivatedNews(t *testing.T) {
	tests := []struct {
		name          string
		access        *helpers.Access
		privilegesDto dto.Privileges
		want          bool
	}{
		{
			name:   "Check privileges [delete_deactivated] for user without scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_active"},
				},
			},
			want: false,
		},
		{
			name:   "Check privileges [delete_deactivated] for user with scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"delete_deactivated"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if got := a.checkCanDeleteDeactivatedNews(tt.privilegesDto); got != tt.want {
				t.Errorf("checkCanDeleteDeactivatedNews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsAccess_CheckCanDeleteNewsOrFail(t *testing.T) {
	type args struct {
		privilegesDto dto.Privileges
		newsEntity    *entity.News
	}
	tests := []struct {
		name    string
		access  *helpers.Access
		args    args
		wantErr bool
	}{
		{
			name:   "Check privileges to delete News [Active=true] for User with scopes [delete_active, delete_deactivated]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active", "delete_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges to delete News [Active=false] for User with scopes [delete_active, delete_deactivated]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active", "delete_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges to delete News [Active=true] for User with scopes [delete_deactivated]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges to delete News [Active=false] for User with scopes [delete_deactivated]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_deactivated"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges to delete News [Active=true] for User with scopes [delete_active]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: false,
		},
		{
			name:   "Check privileges to delete News [Active=false] for User with scopes [delete_active]",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"delete_active"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges to delete News [Active=true] for User without delete scopes",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: &entity.News{
					Active: true,
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges to delete News [Active=false] for User without delete scopes",
			access: helpers.NewAccess(),
			args: args{
				privilegesDto: dto.Privileges{
					Entities: dto.Entities{
						News: []string{"show_active"},
					},
				},
				newsEntity: &entity.News{
					Active: false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if err := a.CheckCanDeleteNewsOrFail(tt.args.privilegesDto, tt.args.newsEntity); (err != nil) != tt.wantErr {
				t.Errorf("CheckCanDeleteNewsOrFail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewsAccess_CheckCanFilterNewsOrFail(t *testing.T) {
	tests := []struct {
		name          string
		access        *helpers.Access
		privilegesDto dto.Privileges
		wantErr       bool
	}{
		{
			name:   "Check privileges [filter] for user without scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"show_active"},
				},
			},
			wantErr: true,
		},
		{
			name:   "Check privileges [filter] for user with scope",
			access: helpers.NewAccess(),
			privilegesDto: dto.Privileges{
				Entities: dto.Entities{
					News: []string{"filter"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewsAccess{
				access: tt.access,
			}
			if err := a.CheckCanFilterNewsOrFail(tt.privilegesDto); (err != nil) != tt.wantErr {
				t.Errorf("CheckCanFilterNewsOrFail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
