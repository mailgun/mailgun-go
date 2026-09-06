Contributing

First off, thank you for considering contributing to Mailgun Go client. Contributions to mailgun-go from the community are welcome and encouraged!

Please create issues for any major changes so it can be discussed first. For small bug fixes feel free to submit the PR directly.

New List endpoints should return `iter.Seq2`(https://pkg.go.dev/iter), e.g.:
```
ListDomainsIter(ctx context.Context, opts *ListDomainsOptions) iter.Seq2[[]mtypes.Domain, error]
```

Make sure to run tests (`make test`) and linters (`make lint`) before creating the PR.

When submitting a pull request, please include the purpose and implementation details.
