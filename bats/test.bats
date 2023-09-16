#!/usr/bin/env bats

assert_status() {
  local expect
  expect="$1"

  [ "${status}" -eq "${expect}" ]
}


@test "should run&show help" {
    run out/bin/fjira --help
    assert_status 0
}
