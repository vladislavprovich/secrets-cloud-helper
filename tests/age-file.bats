#!/usr/bin/env bats

@test "invoke cli - version" {
    run ../dist/go-secretshelper version
    [ "$status" -eq 0 ]
}

@test "invoke cli - run file - nonexisting" {
    run ../dist/go-secretshelper run -c nonex
    [ "$status" -ne 0 ]
}

@test "invoke cli - run file" {
    run ../dist/go-secretshelper run -c ./fixtures/fixture-2.yaml
    [ "$status" -eq 0 ]
    [ -f ./go-secrethelper-test.dat ]
    rm ./go-secrethelper-test.dat
}

@test "invoke cli - run file 3 w/ environment" {
    VAULT_NAME=kv run ../dist/go-secretshelper -e run -c ./fixtures/fixture-3.yaml
    [ "$status" -eq 0 ]
    [ -f ./go-secrethelper-test3.dat ]
    rm ./go-secrethelper-test3.dat
}

@test "invoke cli - run file 4 w/ age transform" {
    age_recipient=age1njkx5t9tcc4gq7c53zzy4sfjq0fscm5uzt5vek5pj2khehcpsfsqwzq9jy run ../dist/go-secretshelper -e run -c ./fixtures/fixture-4.yaml
    [ "$status" -eq 0 ]
    [ -f ./go-secrethelper-test4.dat ]
    FL=$(cat ./go-secrethelper-test4.dat | head -1)
    [ "$FL" = "-----BEGIN AGE ENCRYPTED FILE-----" ]
    rm ./go-secrethelper-test4.dat
}

@test "invoke cli - run file 5 w/ jq transform" {
    run ../dist/go-secretshelper -e run -c ./fixtures/fixture-5.yaml
    [ "$status" -eq 0 ]
    [ -f ./go-secrethelper-test5.dat ]
    FL=$(cat ./go-secrethelper-test5.dat | head -1)
    [ "$FL" = "s3cr3t" ]
    rm ./go-secrethelper-test5.dat
}
