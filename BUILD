load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/jeffbean/drawfh
gazelle(name = "gazelle")

go_library(
    name = "go_default_library",
    srcs = ["helloworld.go"],
    importpath = "github.com/jeffbean/drawfh",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "drawfh",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["helloworld_test.go"],
    embed = [":go_default_library"],
    deps = ["@com_github_stretchr_testify//assert:go_default_library"],
)
