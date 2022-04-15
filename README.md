# s3test

Test S3 stuff

Provide required params through...

### commandline options

    ./s3test --key <yours3key> --secret <yoursecret> --endpoint https://xxx.omgs3stuffs.tld --bucket mahbucket


### environment vars

    export S3TEST_KEY=thisiskey
    export S3TEST_SECRET=verysecret
    export S3TEST_ENDPOINT=https://xxx.omgs3stuffs.tld
    export S3TEST_BUCKET=mahbucket
    ./s3test 

### yaml config file

create a yaml file containing (defaults to ~/.s3test.yaml):

    key: thisiskey
    secret: verysecret
    endpoint: https://xxx.omgs3stuffs.tld
    bucket: mahbucket

then run

    ./s3test --config /some/path/to/s3test.yaml
