load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
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
