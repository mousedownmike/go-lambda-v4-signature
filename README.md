# Go V4 Signature in Lambda

As [described in my post](https://blog.mikedalrymple.com/2019/12/20/v4-signatures-with-go-on-lambda/), this is a simple project that illustrates how to get `Credentials` from a `session.Session` so you can can create a `v4.Signer` and then sign a request all from within an AWS Lambda function.

## Running Locally

This assumes that you AWS credentials configured locally and those credentials permit access to list S3 buckets.

`go run cmd/signed/main.go`

You should see an XML listing of all your S3 buckets.

## Running as a Lambda

Create a zipped binary using the following command from the top of the project directory.

```bash
go build ./...
zip signed.zip signed
```

These are the basic steps to create the function using the console.  If you're here, you probably already know how to do this, there's nothing special about this project.

1. Choose **Create Function**
2. Select **Author from scratch**
   * Function name: `bucketlist`
   * Runtime `Go 1.x`
3. Select **Create Function**
4. Select your `signed.zip` file after selecting the **Upload** button.
5. In the **Handler** field enter `signed`
6. Select the **Save** button in the top right.

The function doesn't expect any input so you can configure a simple test event (e.g. `{}`) and run the function by clicking the **Test** button.  The function will succeed but the output will show that access was denied.  You need to add the `s3:ListAllMyBuckets` policy action to your newly created role.  Run the test again and you'll see the bucket list. 