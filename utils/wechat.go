package utils

import (
	"campus-memory/types/errno"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// WechatConfig 微信小程序配置
type WechatConfig struct {
	AppID  string // 微信小程序 AppID
	Secret string // 微信小程序 Secret
}

// WechatSession 微信登录会话信息
type WechatSession struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识（可选）
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

// GetWechatConfig 从环境变量获取微信配置
func GetWechatConfig() *WechatConfig {
	return &WechatConfig{
		AppID:  os.Getenv("WECHAT_APPID"),
		Secret: os.Getenv("WECHAT_SECRET"),
	}
}

// GetWechatOpenID 调用微信API获取OpenID
func GetWechatOpenID(code string) (*WechatSession, error) {
	// 1. 获取微信配置
	config := GetWechatConfig()

	// 打印调试信息
	fmt.Printf("\n========== 微信登录调试信息 ==========\n")
	fmt.Printf("[1] 接收到的 Code: %s\n", code)
	fmt.Printf("[2] AppID: %s\n", config.AppID)
	fmt.Printf("[3] Secret: %s\n", maskSecret(config.Secret))

	if config.AppID == "" || config.Secret == "" {
		fmt.Printf("[错误] 微信配置未设置\n")
		fmt.Printf("====================================\n\n")
		return nil, errno.New(15003, "微信配置未设置，请联系管理员")
	}
	
	// 2. 构建请求URL
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.AppID,
		config.Secret,
		code,
	)
	fmt.Printf("[4] 请求微信API...\n")

	// 3. 发送HTTP GET请求（设置5秒超时）
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("[错误] 调用微信API失败: %v\n", err)
		fmt.Printf("====================================\n\n")
		return nil, errno.New(15004, "微信服务暂时不可用，请稍后重试")
	}
	defer resp.Body.Close()
	
	// 4. 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[错误] 读取响应失败: %v\n", err)
		fmt.Printf("====================================\n\n")
		return nil, fmt.Errorf("读取微信API响应失败: %w", err)
	}

	fmt.Printf("[5] 微信API响应: %s\n", string(body))

	// 5. 解析JSON响应
	var session WechatSession
	if err := json.Unmarshal(body, &session); err != nil {
		fmt.Printf("[错误] 解析JSON失败: %v\n", err)
		fmt.Printf("====================================\n\n")
		return nil, fmt.Errorf("解析微信API响应失败: %w", err)
	}
	
	// 6. 检查错误码
	if session.ErrCode != 0 {
		fmt.Printf("[错误] 微信API返回错误码: %d, 错误信息: %s\n", session.ErrCode, session.ErrMsg)
		fmt.Printf("====================================\n\n")

		// 根据微信错误码返回具体错误
		switch session.ErrCode {
		case 40029: // code 无效
			return nil, errno.New(15005, "登录凭证已过期，请重新登录")
		case 40163: // code 已使用
			return nil, errno.New(15006, "登录凭证已使用，请重新登录")
		case 40013: // AppID 无效
			return nil, errno.New(15007, "微信配置错误，请联系管理员")
		default:
			return nil, errno.New(15002, fmt.Sprintf("微信登录失败: %s", session.ErrMsg))
		}
	}
	
	// 7. 验证必需字段
	if session.OpenID == "" {
		fmt.Printf("[错误] 微信API未返回OpenID\n")
		fmt.Printf("====================================\n\n")
		return nil, errno.New(15008, "微信登录异常，请稍后重试")
	}

	fmt.Printf("[6] 登录成功! OpenID: %s\n", maskOpenID(session.OpenID))
	fmt.Printf("====================================\n\n")

	return &session, nil
}

// maskSecret 隐藏密钥的部分字符，用于日志输出
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

// maskOpenID 隐藏OpenID的部分字符，用于日志输出
func maskOpenID(openid string) string {
	if len(openid) <= 8 {
		return "****"
	}
	return openid[:4] + "****" + openid[len(openid)-4:]
}
