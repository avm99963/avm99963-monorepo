# Graf archive

This archived copy is an adaptation of the repository found at
https://gerrit.avm99963.com/plugins/gitiles/graf. Specifically, we based this
archive on commit [1082a13f34c15a8cac327b6d1a6e646820b5b375][commit].

## Modifications

The following modifications were made to archive the graf:

- The website is now a static HTML page, instead of rendering the pages via
  PHP.
- The graph data is now saved in a static file, instead of querying it from
  Dario's API (which used to be hosted at
  `https://grafo.dirba.io/api.php?req=getGraph`).
    - The data has also been anonymized, and the login was thus removed as
      well.
- Now the homepage links directly to the graph instead of the login screen.
    - You can still view the login screen anyways by visiting `login.html`. The
      form now always returns shows an error snackbar as if the login failed.
- My email address was censored in the Google Assistant action's privacy
  policy.
- We removed the logic that deletes nodes in the deletable region. This is
  because we based the graf data on the response of the new API hosted at
  https://apps.dafme.upc.edu/graf/, which lacks the nodes that defined this
  region.
- We vendored the import of the Material Design Lite (MDL) library, since the
  CDN is now down and thus the page looked broken (due to the missing
  stylesheet).
    - In fact, this should be upstreamed to [the new fork][fork].

[commit]: https://gerrit.avm99963.com/plugins/gitiles/graf/+/1082a13f34c15a8cac327b6d1a6e646820b5b375
[fork]: https://gitlab.com/fraret/graf
