# goripper
golang soundcloud-ripper

[![Go](https://img.shields.io/badge/Go-v1.21-cyan)]()

a tool whichs basically finds private soundcloud tracks by bruteforcing shareable links of on.soundcloud.com
```
Shareable soundcloud links are shrinked links, who look like this :
=> on.soundcloud.com/XXXXX => code 302 :redirect => soundcloud.com/artist/track),
and this can redirect to public tracks as well as private ones too... bruteforce time.

the script catchs the redirect, and uses a regex to find if the full link contains a private token (public ones does not have any)
and if the regex matches, then a private track has been found.

actually, because of the regex-only filter theres some false-positives:
-old private tracks who kept their same shareable link, which are public now, will stiff have a private token on the link;
-sometimes deleted tracks does match.

this is easy to fix tho, i just need to do it :')
```

---
##### Educationnal purposes only, use it at your own risk.

---
### CLI Arguments
- -h  |  print out the possible commands
---
- -r  |  number of requests
- -t  |  number of threads
---
- -x  |  exports the positive results in a txt file (export.txt), *and updates it if it already exists.*

---
### Usage
you can run the script without any arguments, the default values are enough for testing.

