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
)

type ErrorBuilderFunc func(msg string) error

func DefaultErrorBuilder(msg string) error {
	return errors.New(msg)
}

type RsaTool struct {
	log                  *log.Helper
	PrivateKey           string
	PublicKey            string
	ValidateErrorBuilder ErrorBuilderFunc
	DecryptErrorBuilder  ErrorBuilderFunc
	LengthErrorBuilder   ErrorBuilderFunc
}

func NewRsaTool(logger log.Logger, publicKey string, privateKey string) *RsaTool {
	return &RsaTool{
		log:                  log.NewHelper(logger),
		PrivateKey:           privateKey,
		PublicKey:            publicKey,
		ValidateErrorBuilder: DefaultErrorBuilder,
		DecryptErrorBuilder:  DefaultErrorBuilder,
		LengthErrorBuilder:   DefaultErrorBuilder,
	}
}

func (r *RsaTool) SetErrorBuilder(decrypt ErrorBuilderFunc, validate ErrorBuilderFunc) {
	r.DecryptErrorBuilder = decrypt
	r.ValidateErrorBuilder = validate
}

// ValidateHash 校验 psw 是否与已经过 Hash 的密码 hashed 相等
func (r *RsaTool) ValidateHash(psw string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(psw))
	return err == nil
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
//
// 2. 通过私钥进行解密得到密码明文
//
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

// DecryptAndValidate 将对 pswNeedDecrypt 进行解密得到明文后，与已经生成的 Hash hashedPsw 进行校验，
// 如果两者相等这返回 nil, 否则返回对应的错误信息
func (r *RsaTool) DecryptAndValidate(pswNeedDecrypt, hashedPsw string) error {
	decryptPwd, err := r.DecryptPassword(pswNeedDecrypt)
	if err != nil {
		return err
	}

	if r.ValidateHash(decryptPwd, hashedPsw) {
		return nil
	}

	return r.ValidateErrorBuilder("password error")
}

// DecryptPassword 对密码进行解密，传递的密码已经经过了 Base64 Encode
func (r *RsaTool) DecryptPassword(password string) (string, error) {
	// 解密密码
	decoded, err := base64.StdEncoding.DecodeString(password)
	passwordOutput, err := r.DecryptRSA(decoded, r.PrivateKey)
	passwordOutputStr := string(passwordOutput)

	if len(passwordOutputStr) < 8 {
		return "", r.LengthErrorBuilder("password too short")
	}

	return passwordOutputStr, err
}

// EncryptPassword 对密码进行解密
func (r *RsaTool) EncryptPassword(password string) (string, error) {
	pwdEncrypted, err := r.EncryptRSA([]byte(password), r.PublicKey)
	return string(pwdEncrypted), err
}

func (r *RsaTool) EncryptAndEncodePassword(pwd string) (string, error) {
	pwd, err := r.EncryptPassword(pwd)
	if err != nil {
		return "", err
	}

	decoded := base64.StdEncoding.EncodeToString([]byte(pwd))
	return decoded, nil
}

// DecryptRSA 对数据进行解密操作
func (r *RsaTool) DecryptRSA(src []byte, privKey string) (res []byte, err error) {
	//1.获取秘钥（从本地磁盘读取）
	block, _ := pem.Decode([]byte(privKey))                   //解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes) //还原数据
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	return
}

func (r *RsaTool) EncryptRSA(src []byte, pubKey string) (res []byte, err error) {
	// 2、将得到的字符串解码
	block, _ := pem.Decode([]byte(pubKey))

	// 使用X509将解码之后的数据 解析出来
	//x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes) //对应于生成秘钥的x509.MarshalPKIXPublicKey(&publicKey)
	//keyInit1,err:=x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	//4.使用公钥加密数据
	pKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pKey, src)
	return
}
