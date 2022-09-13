#!/usr/bin/env bats

@test "invoke cli - run file w/ template" {
    run ../dist/go-secretshelper run -c ./fixtures/fixture-template.yaml
    [ "$status" -eq 0 ]
    [ -f ./go-secrethelper-template.dat ]

    CONTENT=`cat ./go-secrethelper-template.dat`
    EXPECTED="sample s3cr3t"
    [ "$CONTENT" == "$EXPECTED" ]
    rm ./go-secrethelper-template.dat

}
