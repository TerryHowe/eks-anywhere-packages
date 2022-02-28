package controllers

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	api "github.com/aws/eks-anywhere-packages/api/v1alpha1"
	"github.com/aws/eks-anywhere-packages/controllers/mocks"
	"github.com/aws/eks-anywhere-packages/pkg/bundle"
	"github.com/aws/eks-anywhere-packages/pkg/testutil"
)

func TestPackageBundleControllerReconcilerReconcile(t *testing.T) {
	t.Parallel()

	discovery := testutil.NewFakeDiscoveryWithDefaults()

	puller := testutil.NewMockPuller()
	bm := bundle.NewBundleManager(log.NullLogger{}, discovery, puller)

	controllerNN := types.NamespacedName{
		Namespace: bundle.ActiveBundleNamespace,
		Name:      bundle.PackageBundleControllerName,
	}
	req := ctrl.Request{
		NamespacedName: controllerNN,
	}

	setMockABC := func(src *api.PackageBundleController) func(ctx context.Context,
		name types.NamespacedName, abc *api.PackageBundleController) error {
		return func(ctx context.Context, name types.NamespacedName,
			target *api.PackageBundleController) error {
			src.DeepCopyInto(target)
			return nil
		}
	}

	setMockBundleList := func(src *api.PackageBundleList) func(ctx context.Context,
		list *api.PackageBundleList, opts *client.ListOptions) error {
		return func(ctx context.Context, target *api.PackageBundleList,
			opts *client.ListOptions) error {
			src.DeepCopyInto(target)
			return nil
		}
	}

	t.Run("marks status inactive if name doesn't match", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mockClient := mocks.NewMockClient(gomock.NewController(t))
		mockABC := &api.PackageBundleController{
			ObjectMeta: metav1.ObjectMeta{
				Name:      bundle.PackageBundleControllerName,
				Namespace: "blah",
			},
			Spec: api.PackageBundleControllerSpec{
				ActiveBundle: "v1.21-1001",
			},
			Status: api.PackageBundleControllerStatus{
				State: api.BundleControllerStateActive,
			},
		}

		inactiveController := types.NamespacedName{
			Name:      "blah",
			Namespace: bundle.ActiveBundleNamespace,
		}
		mockClient.EXPECT().Get(ctx, inactiveController, gomock.Any()).
			DoAndReturn(setMockABC(mockABC))
		mockStatusClient := mocks.NewMockStatusWriter(gomock.NewController(t))
		mockClient.EXPECT().Status().Return(mockStatusClient)
		mockStatusClient.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, abc *api.PackageBundleController,
				opts *client.UpdateOptions) error {
				if abc.Status.State != api.BundleControllerStateIgnored {
					t.Errorf("expected state to be set to Ignored, got %q",
						abc.Status.State)
				}
				return nil
			})

		r := NewPackageBundleControllerReconciler(mockClient, nil, bm,
			log.NullLogger{})
		result, err := r.Reconcile(ctx, ctrl.Request{NamespacedName: inactiveController})
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		if result.Requeue {
			t.Errorf("expected Requeue to be false, got true")
		}
	})

	t.Run("marks status active if name matches", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mockClient := mocks.NewMockClient(gomock.NewController(t))
		mockABC := &api.PackageBundleController{
			ObjectMeta: metav1.ObjectMeta{
				Name:      bundle.PackageBundleControllerName,
				Namespace: bundle.ActiveBundleNamespace,
			},
			Spec: api.PackageBundleControllerSpec{
				ActiveBundle: "v1.21-1001",
			},
			Status: api.PackageBundleControllerStatus{
				State: api.BundleControllerStateIgnored,
			},
		}
		mockBundleList := &api.PackageBundleList{
			Items: []api.PackageBundle{
				*api.MustPackageBundleFromFilename(t, "../api/testdata/bundle_one.yaml"),
			},
		}

		mockClient.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).
			DoAndReturn(setMockABC(mockABC))
		mockStatusClient := mocks.NewMockStatusWriter(gomock.NewController(t))
		mockClient.EXPECT().Status().Return(mockStatusClient)
		mockStatusClient.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, abc *api.PackageBundleController,
				opts *client.UpdateOptions) error {
				if abc.Status.State != api.BundleControllerStateActive {
					t.Errorf("expected state to be set to Active, got %q",
						abc.Status.State)
				}
				return nil
			})
		mockClient.EXPECT().List(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(setMockBundleList(mockBundleList))

		mockPuller := testutil.NewMockPuller()
		mockPuller.WithFileData(t, "../api/testdata/bundle_one.yaml")
		bm := bundle.NewBundleManager(log.NullLogger{}, discovery, mockPuller)
		r := NewPackageBundleControllerReconciler(mockClient, nil, bm,
			log.NullLogger{})
		_, err := r.Reconcile(ctx, req)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
	})

	t.Run("handles deleted package bundle controllers", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		mockClient := mocks.NewMockClient(gomock.NewController(t))

		groupResource := schema.GroupResource{
			Group:    req.Namespace,
			Resource: req.Name,
		}
		notFoundError := errors.NewNotFound(groupResource, req.Name)
		mockClient.EXPECT().Get(ctx, req.NamespacedName, gomock.Any()).
			Return(notFoundError)

		r := NewPackageBundleControllerReconciler(mockClient, nil, bm,
			log.NullLogger{})
		result, err := r.Reconcile(ctx, req)
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		if result.Requeue {
			t.Errorf("expected Requeue to be false, got true")
		}
	})
}