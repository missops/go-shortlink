package utils
import "crypto/sha1"

//ToSha1 ...
func ToSha1(str string) string {
   var (
      sha = sha1.New()
   )
   return string(sha.Sum([]byte(str)))
}