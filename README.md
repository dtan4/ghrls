# ghrls(1)

List & Describe GitHub Releases

```bash
$ ghrls list kubernetes/kubernetes | head
TAG               TYPE           CREATEDAT                        NAME
v1.6.0-alpha.0    TAG
v1.5.3-beta.0     TAG
v1.5.2            TAG+RELEASE    2017-01-12 13:51:15 +0900 JST    v1.5.2
v1.5.2-beta.0     TAG
v1.5.1            TAG+RELEASE    2016-12-14 09:50:36 +0900 JST    v1.5.1
v1.5.1-beta.0     TAG
v1.5.0            TAG+RELEASE    2016-12-13 08:29:43 +0900 JST    v1.5.0
v1.5.0-beta.3     TAG+RELEASE    2016-12-09 06:52:35 +0900 JST    v1.5.0-beta.3
v1.5.0-beta.2     TAG+RELEASE    2016-11-25 07:29:04 +0900 JST    v1.5.0-beta.2

$ ghrls get kubernetes/kubernetes v1.5.2
Tag:         v1.5.2
Commit:      08e099554f3c31f6e6f07b448ab3ed78d0520507
Name:        v1.5.2
Author:      saad-ali
CreatedAt:   2017-01-12 13:51:15 +0900 JST
PublishedAt: 2017-01-12 16:25:50 +0900 JST
URL:         https://github.com/kubernetes/kubernetes/releases/tag/v1.5.2
Assets:      https://github.com/kubernetes/kubernetes/releases/download/v1.5.2/kubernetes.tar.gz

See [kubernetes-announce@](https://groups.google.com/forum/#!forum/kubernetes-announce) and [CHANGELOG](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#v152) for details.

SHA256 for `kubernetes.tar.gz`: `67344958325a70348db5c4e35e59f9c3552232cdc34defb8a0a799ed91c671a3`

Additional binary downloads are linked in the [CHANGELOG](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG.md#downloads-for-v152).
```

## Author

Daisuke Fujita ([@dtan4](https://github.com/dtan4))

## License

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
