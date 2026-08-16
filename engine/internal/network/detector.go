// Package network 探测官方与镜像连通性并判定网络区域。
package network

import "context"

type Probe func(ctx context.Context, url string) error

type Detector struct {
	probe Probe
}

func NewDetector(probe Probe) *Detector {
	return &Detector{probe: probe}
}

type Status struct {
	Region     string `json:"region"` // cn | global
	Reason     string `json:"reason"` // auto | manual
	OfficialOK bool   `json:"official_ok"`
	MirrorOK   bool   `json:"mirror_ok"`
}

// Detect 按 override（settings 里的区域）判定；auto 时官方不可达且镜像可达 → cn。
func (d *Detector) Detect(ctx context.Context, override string) Status {
	if override == "cn" {
		return Status{Region: "cn", Reason: "manual"}
	}
	if override == "global" {
		return Status{Region: "global", Reason: "manual"}
	}
	officialOK := d.probe(ctx, "https://registry.npmjs.org") == nil
	mirrorOK := d.probe(ctx, "https://registry.npmmirror.com") == nil
	if !officialOK && mirrorOK {
		return Status{Region: "cn", Reason: "auto", MirrorOK: true}
	}
	return Status{Region: "global", Reason: "auto", OfficialOK: officialOK, MirrorOK: mirrorOK}
}
