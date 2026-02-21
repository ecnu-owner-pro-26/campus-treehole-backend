package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	if config.AppID == "" || config.Secret == "" {
		return nil, fmt.Errorf("微信配置未设置，请检查环境变量 WECHAT_APPID 和 WECHAT_SECRET")
	}
	
	// 2. 构建请求URL
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.AppID,
		config.Secret,
		code,
	)
	
	// 3. 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("调用微信API失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 4. 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取微信API响应失败: %w", err)
	}
	
	// 5. 解析JSON响应
	var session WechatSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, fmt.Errorf("解析微信API响应失败: %w", err)
	}
	
	// 6. 检查错误码
	if session.ErrCode != 0 {
		return nil, fmt.Errorf("微信API返回错误: [%d] %s", session.ErrCode, session.ErrMsg)
	}
	
	// 7. 验证必需字段
	if session.OpenID == "" {
		return nil, fmt.Errorf("微信API未返回OpenID")
	}
	
	return &session, nil
}
