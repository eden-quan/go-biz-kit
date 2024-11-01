package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type RsaTool struct {
	log            *log.Helper
	PrivateKeyFile string
}

func NewRsaTool(logger log.Logger, privateKeyFile string) *RsaTool {
	return &RsaTool{
		log:            log.NewHelper(logger),
		PrivateKeyFile: privateKeyFile,
	}
}

func (r *RsaTool) ValidateHash(psw string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(psw))
	return err == nil
	//return BcryptValidatePassword(psw, hashed)
}

func (r *RsaTool) GenerateHash(pwd string) (string, error) {
	// DONE: 实现 Crypto::GenerateHash
	//generated := BcryptGenerateHash(pwd)
	//return generated, nil
	generated, err := bcrypt.GenerateFromPassword([]byte(pwd), 4)
	return string(generated), err
}

// GeneratePassword 检验旧的已经加密过后的 password 同时生成新的 Hash, 返回原始的 Password, 新的 Password Hash, 以及是否出错
// 具体步骤为：
//
// 1. 对传递的密码进行 Base64 Decode;
// 2. 通过私钥进行解密得到密码明文
// 3，通过 BCrypt 生成密码的 Hash (一般服务器只存储 Hash)
func (r *RsaTool) GeneratePassword(password string) (string, string, error) {
	originPwd, err := r.DecryptPassword(password)
	if err != nil {
		return "", "", err
	}

	pwd, err := r.GenerateHash(originPwd)
	if err != nil {
		return "", "", err
	}

	return originPwd, pwd, nil
}

// DecryptPassword 对密码进行解密，传递的密码已经经过了 Base64 Encode
func (r *RsaTool) DecryptPassword(password string) (string, error) {
	// 解密密码
	decoded, err := base64.StdEncoding.DecodeString(password)
	passwordOutput, err := r.DecryptRSA(decoded, r.PrivateKeyFile)
	passwordOutputStr := string(passwordOutput)

	if len(passwordOutputStr) < 8 {
		return "", errors.New("password too short")
	}

	// todo: 加密
	return passwordOutputStr, err
}

// EncryptPassword 对密码进行解密
func (r *RsaTool) EncryptPassword(password string) (string, error) {
	pwdEncrypted, err := r.EncryptRSA([]byte(password), r.PrivateKeyFile)
	return string(pwdEncrypted), err
}

// DecryptRSA 对数据进行解密操作
func (r *RsaTool) DecryptRSA(src []byte, path string) (res []byte, err error) {
	//1.获取秘钥（从本地磁盘读取）
	f, err := os.Open(path)
	if err != nil {
		return
	}

	defer func() {
		_ = f.Close()
	}()

	fileInfo, _ := f.Stat()
	b := make([]byte, fileInfo.Size())
	_, _ = f.Read(b)
	block, _ := pem.Decode(b)                                 //解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes) //还原数据
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	return
}

func (r *RsaTool) EncryptRSA(src []byte, path string) (res []byte, err error) {
	//1.获取秘钥（从本地磁盘读取）
	f, err := os.Open(path)
	if err != nil {
		return
	}

	defer func() {
		_ = f.Close()
	}()

	fileInfo, _ := f.Stat()
	b := make([]byte, fileInfo.Size())
	_, _ = f.Read(b)
	// 2、将得到的字符串解码
	block, _ := pem.Decode(b)

	// 使用X509将解码之后的数据 解析出来
	//x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes) //对应于生成秘钥的x509.MarshalPKIXPublicKey(&publicKey)
	//keyInit1,err:=x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	//4.使用公钥加密数据
	pubKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	return
}
