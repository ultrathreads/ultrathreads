package hashid

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/speps/go-hashids/v2"
)

var (
	globalSecret string
	defaultMinLen int
	hashesMap     sync.Map // 用于并发安全地缓存已生成的 HashID 实例
)

// Init 初始化全局配置
func Init(secret string, minLength int) {
	globalSecret = secret
	defaultMinLen = minLength
}

// getHashInstance 内部方法：获取或动态创建指定模型的 HashID 实例
func getHashInstance(modelName string) (*hashids.HashID, error) {
	// 1. 尝试从缓存中读取
	if val, ok := hashesMap.Load(modelName); ok {
		return val.(*hashids.HashID), nil
	}

	// 2. 缓存未命中，创建新实例并写入
	hd := hashids.NewData()
	hd.Salt = fmt.Sprintf("%s_%s", globalSecret, modelName)
	hd.MinLength = defaultMinLen
	
	newHash, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, fmt.Errorf("创建 HashID 实例失败: %w", err)
	}

	// Store 是原子操作，如果多个协程同时创建，只有第一个会被保留，这非常安全
	actual, _ := hashesMap.LoadOrStore(modelName, newHash)
	return actual.(*hashids.HashID), nil
}

// getModelName 内部方法：通过反射安全地提取结构体名称并转小写
func getModelName[T any]() (string, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("仅支持结构体类型，当前传入的是 %v", t.Kind())
	}
	return strings.ToLower(t.Name()), nil
}

// Encode 泛型编码：自动提取 T 的名称作为模型标识
func Encode[T any](realID int64) (string, error) {
	modelName, err := getModelName[T]()
	if err != nil {
		return "", err
	}
	
	h, err := getHashInstance(modelName)
	if err != nil {
		return "", err
	}
	
	return h.EncodeInt64([]int64{realID})
}

// Decode 泛型解码：自动提取 T 的名称作为模型标识
func Decode[T any](hashStr string) (int64, error) {
	modelName, err := getModelName[T]()
	if err != nil {
		return 0, err
	}
	
	h, err := getHashInstance(modelName)
	if err != nil {
		return 0, err
	}
	
	numbers, err := h.DecodeInt64WithError(hashStr)
	if err != nil || len(numbers) == 0 {
		return 0, fmt.Errorf("解码失败或结果为空")
	}
	
	return numbers[0], nil
}