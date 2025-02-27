package linodego

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// LinodeKernel represents a Linode Instance kernel object
type LinodeKernel struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Deprecated   bool   `json:"deprecated"`
	KVM          bool   `json:"kvm"`
	XEN          bool   `json:"xen"`
	PVOPS        bool   `json:"pvops"`
}

// LinodeKernelsPagedResponse represents a Linode kernels API response for listing
type LinodeKernelsPagedResponse struct {
	*PageOptions
	Data []LinodeKernel `json:"data"`
}

func (LinodeKernelsPagedResponse) endpoint(_ ...any) string {
	return "linode/kernels"
}

func (resp *LinodeKernelsPagedResponse) castResult(r *resty.Request, e string) (int, int, error) {
	res, err := coupleAPIErrors(r.SetResult(LinodeKernelsPagedResponse{}).Get(e))
	if err != nil {
		return 0, 0, err
	}
	castedRes := res.Result().(*LinodeKernelsPagedResponse)
	resp.Data = append(resp.Data, castedRes.Data...)
	return castedRes.Pages, castedRes.Results, nil
}

// ListKernels lists linode kernels. This endpoint is cached by default.
func (c *Client) ListKernels(ctx context.Context, opts *ListOptions) ([]LinodeKernel, error) {
	response := LinodeKernelsPagedResponse{}

	if result, err := c.getCachedResponse(response.endpoint()); err != nil {
		return nil, err
	} else if result != nil {
		return result.([]LinodeKernel), nil
	}

	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}

	if err := c.addCachedResponse(response.endpoint(), response.Data, nil); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetKernel gets the kernel with the provided ID. This endpoint is cached by default.
func (c *Client) GetKernel(ctx context.Context, kernelID string) (*LinodeKernel, error) {
	e := fmt.Sprintf("linode/kernels/%s", kernelID)

	if result, err := c.getCachedResponse(e); err != nil {
		return nil, err
	} else if result != nil {
		result := result.(LinodeKernel)
		return &result, nil
	}

	req := c.R(ctx).SetResult(&LinodeKernel{})
	r, err := coupleAPIErrors(req.Get(e))
	if err != nil {
		return nil, err
	}

	if err := c.addCachedResponse(e, r.Result(), nil); err != nil {
		return nil, err
	}

	return r.Result().(*LinodeKernel), nil
}
