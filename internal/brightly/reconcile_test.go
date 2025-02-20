package brightly

import (
	"fmt"
	"github.com/launchdarkly/go-server-sdk-evaluation/v3/ldmodel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Reconcile(t *testing.T) {
	tests := []struct {
		name     string
		old      relayArchiveBuilder
		new      relayArchiveBuilder
		expected relayArchiveBuilder
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "no change",
			old: relayArchive().
				env(env("staging").version(1).dataId("2").
					flag(booleanFlag("flag1").variation(false).version(1))),
			new: relayArchive().
				env(env("staging").
					flag(booleanFlag("flag1").variation(false))),
			expected: relayArchive().
				env(env("staging").version(1).dataId("2").
					flag(booleanFlag("flag1").variation(false).version(1))),
			wantErr: assert.NoError,
		},
		{
			name: "toggle flag",
			old: relayArchive().
				env(env("staging").version(1).dataId("2").
					flag(booleanFlag("flag1").variation(false).version(1))),
			new: relayArchive().
				env(env("staging").
					flag(booleanFlag("flag1").variation(true))),
			expected: relayArchive().
				env(env("staging").version(1).dataId("3").
					flag(booleanFlag("flag1").variation(true).version(2))),
			wantErr: assert.NoError,
		},
		{
			name: "new flag",
			old: relayArchive().
				env(env("staging").version(1).dataId("2").
					flag(booleanFlag("flag1").variation(false).version(1))),
			new: relayArchive().
				env(env("staging").
					flag(booleanFlag("flag1").variation(false)).
					flag(booleanFlag("flag2").variation(true)),
				),
			expected: relayArchive().
				env(env("staging").version(1).dataId("3").
					flag(booleanFlag("flag1").variation(false).version(1)).
					flag(booleanFlag("flag2").variation(true).version(1)),
				),
			wantErr: assert.NoError,
		},
		{
			name: "deleted flag",
			old: relayArchive().
				env(env("staging").version(1).dataId("2").
					flag(booleanFlag("flag1").variation(false).version(1))),
			new: relayArchive().
				env(env("staging").version(1)),
			expected: relayArchive().
				env(env("staging").version(1).dataId("3").
					flag(booleanFlag("flag1").variation(false).version(2).deleted(true))),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reconcileArchives(tt.old.RelayArchive, tt.new.RelayArchive)
			if !tt.wantErr(t, err, fmt.Sprintf("reconcileArchives(%v, %v)", tt.old.RelayArchive, tt.new.RelayArchive)) {
				return
			}
			assert.Equalf(t, tt.expected.RelayArchive, got, "reconcileArchives(%v, %v)", tt.old.RelayArchive, tt.new.RelayArchive)
		})
	}
}

type relayArchiveBuilder struct {
	RelayArchive
}

func relayArchive() relayArchiveBuilder {
	return relayArchiveBuilder{
		RelayArchive: RelayArchive{map[string]Env{}},
	}
}

func (b relayArchiveBuilder) env(envBuilder envBuilder) relayArchiveBuilder {
	b.envs[envBuilder.Env.metadata.EnvMetadata.EnvKey] = envBuilder.Env
	return b
}

type envBuilder struct {
	Env
}

func env(key string) envBuilder {
	archiveEnv := RelayArchiveEnv{
		EnvMetadata: RelayArchiveEnvMetadata{
			EnvID:   key,
			EnvKey:  key,
			EnvName: key,
		},
	}

	return envBuilder{Env{
		archiveEnv,
		RelayArchiveData{
			Flags: map[string]ldmodel.FeatureFlag{},
		}}}
}
func (b envBuilder) dataId(dataId string) envBuilder {
	b.Env.metadata.DataId = dataId
	return b
}

func (b envBuilder) version(version int) envBuilder {
	b.Env.metadata.EnvMetadata.Version = version
	return b
}

func (b envBuilder) flag(flagBuilder flagBuilder) envBuilder {
	flag := flagBuilder.toLdFlag()
	b.Env.data.Flags[flag.Key] = flag
	return b
}

type flagBuilder interface {
	toLdFlag() ldmodel.FeatureFlag
}

type flagBuilderBase struct {
	FlagBase
	versionV  int
	isDeleted bool
}

type booleanFlagBuilder struct {
	flagBuilderBase
	variationV bool
}

func booleanFlag(key string) booleanFlagBuilder {
	return booleanFlagBuilder{
		flagBuilderBase: flagBuilderBase{FlagBase: FlagBase{key: key}},
	}
}

func (b booleanFlagBuilder) deleted(isDeleted bool) booleanFlagBuilder {
	b.isDeleted = isDeleted
	return b
}

func (b booleanFlagBuilder) variation(variation bool) booleanFlagBuilder {
	b.variationV = variation
	return b
}

func (b booleanFlagBuilder) version(version int) booleanFlagBuilder {
	b.versionV = version
	return b
}

func (b booleanFlagBuilder) toLdFlag() ldmodel.FeatureFlag {
	f := FlagBoolean{
		FlagBase:  b.FlagBase,
		Variation: b.variationV,
	}

	ldFlag := f.ToLdFlag()
	ldFlag.Version = b.versionV
	ldFlag.Deleted = b.isDeleted

	return ldFlag
}

func Test_compareMaps(t *testing.T) {
	type testCase struct {
		name     string
		old      map[string]string
		new      map[string]string
		expected compareResult
	}
	tests := []testCase{
		{
			name: "nil input",
			old:  nil,
			new:  nil,
			expected: compareResult{
				new:      []string{},
				existing: []string{},
				deleted:  []string{},
			},
		},
		{
			name: "nil new",
			old:  map[string]string{"aKey": "aValue", "bKey": "bValue"},
			new:  nil,
			expected: compareResult{
				new:      []string{},
				existing: []string{},
				deleted:  []string{"aKey", "bKey"},
			},
		},
		{
			name: "nil old",
			old:  nil,
			new:  map[string]string{"aKey": "aValue", "bKey": "bValue"},
			expected: compareResult{
				new:      []string{"aKey", "bKey"},
				existing: []string{},
				deleted:  []string{},
			},
		},
		{
			name: "same",
			old:  map[string]string{"aKey": "aValue", "bKey": "bValue"},
			new:  map[string]string{"aKey": "aValue", "bKey": "bValue"},
			expected: compareResult{
				new:      []string{},
				existing: []string{"aKey", "bKey"},
				deleted:  []string{},
			},
		},
		{
			name: "mixed",
			old:  map[string]string{"deletedKey": "deletedValue", "bKey": "bValue"},
			new:  map[string]string{"newKey": "newValue", "bKey": "bValue"},
			expected: compareResult{
				new:      []string{"newKey"},
				existing: []string{"bKey"},
				deleted:  []string{"deletedKey"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := compareMapKeys(tt.old, tt.new)
			assert.ElementsMatch(t, tt.expected.new, actual.new, "compareMapKeys(%v, %v)", tt.old, tt.new)
			assert.ElementsMatch(t, tt.expected.existing, actual.existing, "compareMapKeys(%v, %v)", tt.old, tt.new)
			assert.ElementsMatch(t, tt.expected.deleted, actual.deleted, "compareMapKeys(%v, %v)", tt.old, tt.new)
		})
	}
}
