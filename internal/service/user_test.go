package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
	"tonotdolist/common"
	"tonotdolist/internal/model"
	mock_clock "tonotdolist/mocks/clock"
	mock_util "tonotdolist/mocks/id"
	mock_repository "tonotdolist/mocks/repository"
	"tonotdolist/pkg/clock"
)

const (
	bcryptCost       = bcrypt.MinCost
	clockTime  int64 = 1000
	id               = "5627a05a-111a-4ba8-977b-9ba0d17c4793"

	sessionLength int64 = 1000
	expectedExp         = clockTime + sessionLength

	email = "t@test.t"
)

var anyError = errors.New("any")

func conf() *viper.Viper {
	c := viper.New()
	c.Set("auth.bcryptCost", bcryptCost) // tested in register success test case under mock user repo
	c.Set("auth.sessionLength", sessionLength)

	return c
}

func init() {
	// tested in register & login test cases @ the end result check
}

func TestUserService_GetSession(t *testing.T) {
	type testcase struct {
		tcName      string
		sessionId   string
		userId      string
		expectedErr error
		validResult bool
		configure   func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository)
	}

	cases := []testcase{
		{
			tcName:      "Success",
			userId:      "very epic userid",
			sessionId:   "super good session id",
			expectedErr: nil,
			validResult: true,
			configure: func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockSessionRepo.EXPECT().GetSession(gomock.Any(), gomock.Eq(tc.userId)).Return(&common.UserSession{UserID: tc.userId, Expire: clockTime}, nil) // last second
			},
		},
		{
			tcName:      "ExpiredSession",
			userId:      "very epic userid",
			sessionId:   "super good session id",
			expectedErr: common.ErrUnauthorized,
			validResult: false,
			configure: func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockSessionRepo.EXPECT().GetSession(gomock.Any(), gomock.Eq(tc.userId)).Return(&common.UserSession{UserID: tc.userId, Expire: clockTime - 1}, nil)
			},
		},
		{
			tcName:      "SessionNotFoundInRepo",
			userId:      "very epic userid",
			sessionId:   "super good session id",
			expectedErr: common.ErrUnauthorized,
			validResult: false,
			configure: func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockSessionRepo.EXPECT().GetSession(gomock.Any(), gomock.Eq(tc.userId)).Return(nil, common.ErrNotFound)
			},
		},
		{
			tcName:      "SessionRepoFail",
			userId:      "very epic userid",
			sessionId:   "super good session id",
			expectedErr: errors.New("sad repo error"),
			validResult: false,
			configure: func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockSessionRepo.EXPECT().GetSession(gomock.Any(), gomock.Eq(tc.userId)).Return(nil, tc.expectedErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tcName, func(t *testing.T) {
			t.Parallel()

			config := conf()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
			mockSessionRepo := mock_repository.NewMockSessionRepository(ctrl)
			mockIdProvider := mock_util.NewMockIDProvider(ctrl)
			mockIdProvider.EXPECT().NewID().Return(id, nil).Times(0)

			tc.configure(t, &tc, mockSessionRepo)

			userService := NewUserService(getClock(ctrl), mockIdProvider, mockUserRepo, mockSessionRepo, config)
			userId, err := userService.GetSession(context.Background(), tc.userId)

			if errors.Is(tc.expectedErr, anyError) {
				assert.NotNil(t, err, "Expected error. ")
			} else {
				assert.ErrorIs(t, err, tc.expectedErr, "Error does not match expected error.")
			}

			if tc.validResult {
				assert.Equal(t, tc.userId, userId, "Test case user id does not match with result user id.")
			}
		})
	}
}

func TestUserService_Register(t *testing.T) {
	type testcase struct {
		tcName        string
		email         string
		password      string
		expectedErr   error
		customIdLogic bool
		validResult   bool
		configure     func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider)
	}

	cases := []testcase{
		{
			tcName:        "Success",
			email:         email,
			password:      "Password1_",
			expectedErr:   nil,
			customIdLogic: false,
			validResult:   true,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
				mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, actualUser *model.User) error {
						assert.Equal(t, tc.email, actualUser.Email)
						assert.Equal(t, id, actualUser.UserId)
						err := bcrypt.CompareHashAndPassword([]byte(actualUser.Password), []byte(tc.password))
						assert.NoError(t, err)

						// verify bcrypt cost
						cost, err := bcrypt.Cost([]byte(actualUser.Password))
						require.NoError(t, err, "Unexpected error while trying to get cost of generated password.")
						assert.Equal(t, bcryptCost, cost, "Password cost does not match expected cost. ")

						return nil
					}).Times(1)

				mockSessionRepo.EXPECT().AddSession(gomock.Any(), gomock.Eq(id), gomock.Eq(id), expectedExp)
			},
		},
		{
			tcName:        "Conflict",
			email:         email,
			password:      "Password1_",
			expectedErr:   common.ErrConflict,
			customIdLogic: false,
			validResult:   false,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
				mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(common.ErrConflict).Times(1)
			},
		},
		{
			tcName:        "BadPassword",
			email:         email,
			password:      "                 ",
			expectedErr:   common.ErrBadPassword,
			customIdLogic: false,
			validResult:   false,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
			},
		},
		{
			tcName:        "ShortPassword",
			email:         email,
			password:      "Aa1_",
			expectedErr:   common.ErrPasswordTooShort,
			customIdLogic: false,
			validResult:   false,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
			},
		},
		{
			tcName:        "LongPassword",
			email:         email,
			password:      "Wf9$kp!7Lt#vD@zXbQ1uSe&3Mo*HjArTgN0YwIcVx+ElRsK8J-ZP4yFUnCmOq5B6dAAkdkkdf",
			expectedErr:   common.ErrPasswordTooLong,
			customIdLogic: false,
			validResult:   false,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
			},
		},
		{
			tcName:        "SessionRepoFail",
			email:         email,
			password:      "Wf9$kp!7Lt#vD@zXbQ1uSe&3Mo*HjArTgN0YwIcVx+-ZP4yFUnCkkdf",
			expectedErr:   errors.New("session add fail"),
			customIdLogic: false,
			validResult:   false,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
				mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mockSessionRepo.EXPECT().AddSession(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.expectedErr)
			},
		},
		{
			tcName:        "CreateUserIdProviderFail",
			email:         email,
			password:      "Wf9$kp!7Lt#vD@zXbQ1uSe&3Mo*HjArTgN0YwIcVx+-ZP4yFUnCkkdf",
			expectedErr:   errors.New("id provide fail"),
			customIdLogic: true,
			validResult:   false,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
				mockIdProvider.EXPECT().NewID().Return("best id ever", tc.expectedErr).AnyTimes()
			},
		},
		{
			tcName:        "CreateSessionIdProviderFail",
			email:         email,
			password:      "Wf9$kp!7Lt#vD@zXbQ1uSe&3Mo*HjArTgN0YwIcVx+-ZP4yFUnCkkdf",
			expectedErr:   errors.New("id provide fail"),
			customIdLogic: true,
			validResult:   false,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
				mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any())
				mockIdProvider.EXPECT().NewID().Return("good here", nil).Times(1)
				mockIdProvider.EXPECT().NewID().Return("oopsie", tc.expectedErr).Times(1)
			},
		},
		{
			tcName:        "HashPasswordFail",
			email:         email,
			password:      "Wf9$kp!7Lt#vD@zXbQ1uSe&3Mo*HjArTgN0YwIcVx+-ZP4yFUnCkkdf",
			expectedErr:   anyError,
			customIdLogic: false,
			validResult:   false,
			configure: func(t *testing.T, tc *testcase, config *viper.Viper, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository, mockIdProvider *mock_util.MockIDProvider) {
				config.Set(bcryptCostKey, 999) // exceeding the max cost will return error
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tcName, func(t *testing.T) {
			t.Parallel()

			config := conf()

			request := &common.UserRegisterRequest{
				Email:    tc.email,
				Password: tc.password,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
			mockSessionRepo := mock_repository.NewMockSessionRepository(ctrl)
			mockIdProvider := mock_util.NewMockIDProvider(ctrl)
			if !tc.customIdLogic {
				mockIdProvider.EXPECT().NewID().Return(id, nil).AnyTimes()
			}

			tc.configure(t, &tc, config, mockUserRepo, mockSessionRepo, mockIdProvider)

			userService := NewUserService(getClock(ctrl), mockIdProvider, mockUserRepo, mockSessionRepo, config)
			sid, err := userService.Register(context.Background(), request)

			if errors.Is(tc.expectedErr, anyError) {
				assert.NotNil(t, err, "Error expected.")
			} else {
				assert.ErrorIs(t, err, tc.expectedErr, "Error does not match expected error.")
			}

			if tc.validResult {
				assert.Equal(t, sid, id, "Expected final session id to match the correct session id.")
			} else {
				assert.NotEqual(t, sid, id, "Expected final session id to not match the correct session id.")
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	const (
		password = "word"
	)

	type testcase struct {
		tcName      string
		email       string
		password    string
		expectedErr error
		validResult bool
		configure   func(t *testing.T, tc *testcase, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository)
	}

	cases := []testcase{
		{
			tcName:      "Success",
			email:       email,
			password:    password,
			expectedErr: nil,
			validResult: true,
			configure: func(t *testing.T, tc *testcase, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockUserRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(tc.email)).DoAndReturn(func(ctx context.Context, argEmail string) (*model.User, error) {
					hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

					require.NoError(t, err, "error generating hashed password")

					return &model.User{
						UserId:   id,
						Password: string(hashedPwd),
					}, nil
				})

				mockSessionRepo.EXPECT().AddSession(gomock.Any(), gomock.Eq(id), gomock.Eq(id), expectedExp)
			},
		},
		{
			tcName:      "WrongPassword",
			email:       email,
			password:    password,
			expectedErr: common.ErrUnauthorized,
			validResult: false,
			configure: func(t *testing.T, tc *testcase, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockUserRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(tc.email)).DoAndReturn(func(ctx context.Context, argEmail string) (*model.User, error) {
					return &model.User{
						UserId:   id,
						Password: "securehash2000",
					}, nil
				})
			},
		},
		{
			tcName:      "UserRepoFetchFail",
			email:       email,
			password:    password,
			expectedErr: errors.New("funny error"),
			validResult: false,
			configure: func(t *testing.T, tc *testcase, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockUserRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(tc.email)).DoAndReturn(func(ctx context.Context, argEmail string) (*model.User, error) {
					return nil, tc.expectedErr
				})
			},
		},
		{
			tcName:      "SessionCreateFail",
			email:       email,
			password:    password,
			expectedErr: errors.New("funny error"),
			validResult: false,
			configure: func(t *testing.T, tc *testcase, mockUserRepo *mock_repository.MockUserRepository, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockUserRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(tc.email)).DoAndReturn(func(ctx context.Context, argEmail string) (*model.User, error) {
					hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

					require.NoErrorf(t, err, "error generating hashed password")

					return &model.User{
						UserId:   id,
						Password: string(hashedPwd),
					}, nil
				})

				mockSessionRepo.EXPECT().AddSession(gomock.Any(), gomock.Eq(id), gomock.Eq(id), expectedExp).Return(tc.expectedErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tcName, func(t *testing.T) {
			t.Parallel()

			config := conf()

			request := &common.UserLoginRequest{
				Email:    tc.email,
				Password: tc.password,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
			mockSessionRepo := mock_repository.NewMockSessionRepository(ctrl)
			mockIdProvider := mock_util.NewMockIDProvider(ctrl)
			mockIdProvider.EXPECT().NewID().Return(id, nil).AnyTimes()

			tc.configure(t, &tc, mockUserRepo, mockSessionRepo)

			userService := NewUserService(getClock(ctrl), mockIdProvider, mockUserRepo, mockSessionRepo, config)
			sid, err := userService.Login(context.Background(), request)

			if errors.Is(tc.expectedErr, anyError) {
				assert.NotNil(t, err, "Expected error. ")
			} else {
				assert.ErrorIs(t, err, tc.expectedErr, "Error does not match expected error.")
			}

			if tc.validResult {
				assert.Equal(t, sid, id, "Expected session id to match the correct session id.")
			}
		})
	}
}

func TestUserService_Logout(t *testing.T) {
	type testcase struct {
		tcName      string
		sessionId   string
		userId      string
		expectedErr error
		configure   func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository)
	}

	cases := []testcase{
		{
			tcName:      "Success",
			sessionId:   "amazing sid",
			userId:      "very epic userid",
			expectedErr: nil,
			configure: func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockSessionRepo.EXPECT().DeleteSession(gomock.Any(), gomock.Eq(tc.sessionId), gomock.Eq(tc.userId)).Return(nil)
			},
		},
		{
			tcName:      "RepoFail",
			sessionId:   "amazing sid",
			userId:      "very epic userid",
			expectedErr: errors.New("session repo error"),
			configure: func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockSessionRepo.EXPECT().DeleteSession(gomock.Any(), gomock.Eq(tc.sessionId), gomock.Eq(tc.userId)).Return(tc.expectedErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tcName, func(t *testing.T) {
			t.Parallel()

			config := conf()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
			mockSessionRepo := mock_repository.NewMockSessionRepository(ctrl)
			mockIdProvider := mock_util.NewMockIDProvider(ctrl)
			mockIdProvider.EXPECT().NewID().Return(id, nil).Times(0)

			tc.configure(t, &tc, mockSessionRepo)

			userService := NewUserService(getClock(ctrl), mockIdProvider, mockUserRepo, mockSessionRepo, config)
			err := userService.Logout(context.Background(), tc.sessionId, tc.userId)

			if errors.Is(tc.expectedErr, anyError) {
				assert.NotNil(t, err, "Expected error. ")
			} else {
				assert.ErrorIs(t, err, tc.expectedErr, "Error does not match expected error.")
			}
		})
	}
}

func TestUserService_LogoutAll(t *testing.T) {
	type testcase struct {
		tcName      string
		userId      string
		expectedErr error
		configure   func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository)
	}

	cases := []testcase{
		{
			tcName:      "Success",
			userId:      "very epic userid",
			expectedErr: nil,
			configure: func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockSessionRepo.EXPECT().DeleteAllUserSession(gomock.Any(), gomock.Eq(tc.userId)).Return(nil)
			},
		},
		{
			tcName:      "RepoFail",
			userId:      "very epic userid",
			expectedErr: errors.New("session repo error"),
			configure: func(t *testing.T, tc *testcase, mockSessionRepo *mock_repository.MockSessionRepository) {
				mockSessionRepo.EXPECT().DeleteAllUserSession(gomock.Any(), gomock.Eq(tc.userId)).Return(tc.expectedErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tcName, func(t *testing.T) {
			t.Parallel()

			config := conf()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
			mockSessionRepo := mock_repository.NewMockSessionRepository(ctrl)
			mockIdProvider := mock_util.NewMockIDProvider(ctrl)
			mockIdProvider.EXPECT().NewID().Return(id, nil).Times(0)

			tc.configure(t, &tc, mockSessionRepo)

			userService := NewUserService(getClock(ctrl), mockIdProvider, mockUserRepo, mockSessionRepo, config)
			err := userService.LogoutAll(context.Background(), tc.userId)

			if errors.Is(tc.expectedErr, anyError) {
				assert.NotNil(t, err, "Expected error. ")
			} else {
				assert.ErrorIs(t, err, tc.expectedErr, "Error does not match expected error.")
			}
		})
	}
}

func getClock(ctrl *gomock.Controller) clock.Clock {
	mockClock := mock_clock.NewMockClock(ctrl)
	mockClock.EXPECT().Now().Return(time.Unix(clockTime, 0)).AnyTimes()

	return mockClock
}
