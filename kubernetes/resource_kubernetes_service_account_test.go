package kubernetes

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAccKubernetesServiceAccount_basic(t *testing.T) {
	var conf api.ServiceAccount
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	resourceName := "kubernetes_service_account.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceAccountConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceAccountExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.annotations.TestAnnotationTwo", "two"),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.labels.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.labels.TestLabelTwo", "two"),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.labels.TestLabelThree", "three"),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.name", name),
					resource.TestCheckResourceAttrSet(resourceName, "metadata.0.generation"),
					resource.TestCheckResourceAttrSet(resourceName, "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "metadata.0.uid"),
					resource.TestCheckResourceAttr(resourceName, "secret.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "image_pull_secret.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "automount_service_account_token", "true"),
					testAccCheckServiceAccountImagePullSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-three$"),
						regexp.MustCompile("^" + name + "-four$"),
					}),
					testAccCheckServiceAccountSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-one$"),
						regexp.MustCompile("^" + name + "-two$"),
						regexp.MustCompile("^" + name + "-token-[a-z0-9]+$"),
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "automount_service_account_token"},
			},
		},
	})
}

func TestAccKubernetesServiceAccount_default_secret(t *testing.T) {
	var conf api.ServiceAccount
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	resourceName := "kubernetes_service_account_v1.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfClusterVersionGreaterThanOrEqual(t, "1.24.0")
		},
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceAccountConfig_default_secret(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceAccountExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "metadata.0.name", name),
					resource.TestCheckResourceAttrSet(resourceName, "default_secret_name"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "automount_service_account_token"},
			},
		},
	})
}

func TestAccKubernetesServiceAccount_automount(t *testing.T) {
	var conf api.ServiceAccount
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     "kubernetes_service_account.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceAccountConfig_automount(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceAccountExists("kubernetes_service_account.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.TestAnnotationTwo", "two"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.TestLabelTwo", "two"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.TestLabelThree", "three"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "secret.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "image_pull_secret.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "automount_service_account_token", "false"),
					testAccCheckServiceAccountImagePullSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-three$"),
						regexp.MustCompile("^" + name + "-four$"),
					}),
					testAccCheckServiceAccountSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-one$"),
						regexp.MustCompile("^" + name + "-two$"),
						regexp.MustCompile("^" + name + "-token-[a-z0-9]+$"),
					}),
				),
			},
		},
	})
}

func TestAccKubernetesServiceAccount_update(t *testing.T) {
	var conf api.ServiceAccount
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     "kubernetes_service_account.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceAccountConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceAccountExists("kubernetes_service_account.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.TestAnnotationTwo", "two"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.TestLabelTwo", "two"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.TestLabelThree", "three"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "secret.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "image_pull_secret.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "automount_service_account_token", "true"),
					testAccCheckServiceAccountImagePullSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-three$"),
						regexp.MustCompile("^" + name + "-four$"),
					}),
					testAccCheckServiceAccountSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-one$"),
						regexp.MustCompile("^" + name + "-two$"),
						regexp.MustCompile("^" + name + "-token-[a-z0-9]+$"),
					}),
				),
			},
			{
				Config: testAccKubernetesServiceAccountConfig_modified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceAccountExists("kubernetes_service_account.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.Different", "1234"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.TestLabelThree", "three"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "secret.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "image_pull_secret.#", "3"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "automount_service_account_token", "false"),
					testAccCheckServiceAccountImagePullSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-three$"),
						regexp.MustCompile("^" + name + "-four$"),
					}),
					testAccCheckServiceAccountSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-one$"),
						regexp.MustCompile("^" + name + "-two$"),
						regexp.MustCompile("^" + name + "-token-[a-z0-9]+$"),
					}),
				),
			},
			{
				Config: testAccKubernetesServiceAccountConfig_noAttributes(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceAccountExists("kubernetes_service_account.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.%", "0"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.%", "0"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "secret.#", "0"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "image_pull_secret.#", "0"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "automount_service_account_token", "true"),
					testAccCheckServiceAccountImagePullSecrets(&conf, []*regexp.Regexp{}),
					testAccCheckServiceAccountSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + name + "-token-[a-z0-9]+$"),
					}),
				),
			},
		},
	})
}

func TestAccKubernetesServiceAccount_generatedName(t *testing.T) {
	var conf api.ServiceAccount
	prefix := "tf-acc-test-gen-"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     "kubernetes_service_account.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceAccountConfig_generatedName(prefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceAccountExists("kubernetes_service_account.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.annotations.%", "0"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.labels.%", "0"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "metadata.0.generate_name", prefix),
					resource.TestMatchResourceAttr("kubernetes_service_account.test", "metadata.0.name", regexp.MustCompile("^"+prefix)),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service_account.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service_account.test", "automount_service_account_token", "true"),
					testAccCheckServiceAccountImagePullSecrets(&conf, []*regexp.Regexp{}),
					testAccCheckServiceAccountSecrets(&conf, []*regexp.Regexp{
						regexp.MustCompile("^" + prefix + "[a-z0-9]+-token-[a-z0-9]+$"),
					}),
				),
			},
		},
	})
}

func testAccCheckServiceAccountImagePullSecrets(m *api.ServiceAccount, expected []*regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(expected) == 0 && len(m.ImagePullSecrets) == 0 {
			return nil
		}

		if !matchLocalObjectReferenceName(m.ImagePullSecrets, expected) {
			return fmt.Errorf("%s image pull secrets don't match.\nExpected: %q\nGiven: %q",
				m.Name, expected, m.ImagePullSecrets)
		}

		return nil
	}
}

func matchLocalObjectReferenceName(lor []api.LocalObjectReference, expected []*regexp.Regexp) bool {
	for _, r := range expected {
		for _, ps := range lor {
			matched := r.MatchString(ps.Name)
			if matched {
				return true
			}
		}
	}
	return false
}

func testAccCheckServiceAccountSecrets(m *api.ServiceAccount, expected []*regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if clusterVersionGreaterThanOrEqual("1.24.0") {
			return nil
		}
		if len(expected) == 0 && len(m.Secrets) == 0 {
			return nil
		}
		if !matchObjectReferenceName(m.Secrets, expected) {
			return fmt.Errorf("%s secrets don't match.\nExpected: %q\nGiven: %q",
				m.Name, expected, m.Secrets)
		}
		return nil
	}
}

func matchObjectReferenceName(lor []api.ObjectReference, expected []*regexp.Regexp) bool {
	for _, r := range expected {
		for _, ps := range lor {
			matched := r.MatchString(ps.Name)
			if matched {
				return true
			}
		}
	}
	return false
}

func testAccCheckKubernetesServiceAccountDestroy(s *terraform.State) error {
	conn, err := testAccProvider.Meta().(KubeClientsets).MainClientset()

	if err != nil {
		return err
	}
	ctx := context.TODO()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_service_account" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Service Account still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKubernetesServiceAccountExists(n string, obj *api.ServiceAccount) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn, err := testAccProvider.Meta().(KubeClientsets).MainClientset()
		if err != nil {
			return err
		}
		ctx := context.TODO()

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		out, err := conn.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesServiceAccountConfig_basic(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service_account" "test" {
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      TestAnnotationTwo = "two"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  secret {
    name = "${kubernetes_secret.one.metadata.0.name}"
  }

  secret {
    name = "${kubernetes_secret.two.metadata.0.name}"
  }

  image_pull_secret {
    name = "${kubernetes_secret.three.metadata.0.name}"
  }

  image_pull_secret {
    name = "${kubernetes_secret.four.metadata.0.name}"
  }
}

resource "kubernetes_secret" "one" {
  metadata {
    name = "%s-one"
  }
}

resource "kubernetes_secret" "two" {
  metadata {
    name = "%s-two"
  }
}

resource "kubernetes_secret" "three" {
  metadata {
    name = "%s-three"
  }
}

resource "kubernetes_secret" "four" {
  metadata {
    name = "%s-four"
  }
}
`, name, name, name, name, name)
}

func testAccKubernetesServiceAccountConfig_default_secret(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service_account_v1" "test" {
  metadata {
    name = "%s"
  }
}`, name)
}

func testAccKubernetesServiceAccountConfig_modified(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service_account" "test" {
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      Different         = "1234"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  secret {
    name = "${kubernetes_secret.one.metadata.0.name}"
  }

  image_pull_secret {
    name = "${kubernetes_secret.two.metadata.0.name}"
  }

  image_pull_secret {
    name = "${kubernetes_secret.three.metadata.0.name}"
  }

  image_pull_secret {
    name = "${kubernetes_secret.four.metadata.0.name}"
  }

  automount_service_account_token = false
}

resource "kubernetes_secret" "one" {
  metadata {
    name = "%s-one"
  }
}

resource "kubernetes_secret" "two" {
  metadata {
    name = "%s-two"
  }
}

resource "kubernetes_secret" "three" {
  metadata {
    name = "%s-three"
  }
}

resource "kubernetes_secret" "four" {
  metadata {
    name = "%s-four"
  }
}
`, name, name, name, name, name)
}

func testAccKubernetesServiceAccountConfig_noAttributes(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service_account" "test" {
  metadata {
    name = "%s"
  }
}
`, name)
}

func testAccKubernetesServiceAccountConfig_generatedName(prefix string) string {
	return fmt.Sprintf(`resource "kubernetes_service_account" "test" {
  metadata {
    generate_name = "%s"
  }
}
`, prefix)
}

func testAccKubernetesServiceAccountConfig_automount(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service_account" "test" {
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      TestAnnotationTwo = "two"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  secret {
    name = "${kubernetes_secret.one.metadata.0.name}"
  }

  secret {
    name = "${kubernetes_secret.two.metadata.0.name}"
  }

  image_pull_secret {
    name = "${kubernetes_secret.three.metadata.0.name}"
  }

  image_pull_secret {
    name = "${kubernetes_secret.four.metadata.0.name}"
  }

  automount_service_account_token = false
}

resource "kubernetes_secret" "one" {
  metadata {
    name = "%s-one"
  }
}

resource "kubernetes_secret" "two" {
  metadata {
    name = "%s-two"
  }
}

resource "kubernetes_secret" "three" {
  metadata {
    name = "%s-three"
  }
}

resource "kubernetes_secret" "four" {
  metadata {
    name = "%s-four"
  }
}
`, name, name, name, name, name)
}
