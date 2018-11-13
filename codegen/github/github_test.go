package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, tagsListResponse)
	}))
	defer ts.Close()

	githubClient := NewClient()
	githubClient.baseURL = ts.URL + "/"
	env := &testEnv{githubClient}

	t.Run("TestListTags", env.TestListTags)
	t.Run("TestLatestTag", env.TestLatestTag)
}

type testEnv struct {
	githubClient *Client
}

func (e *testEnv) TestListTags(t *testing.T) {
	tags, err := e.githubClient.ListTags("hashicorp", "vault")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 30 {
		t.Fatal("expected 30 tags")
	}
	// Spot check a SHA.
	if tags[1].Commit.SHA != "612120e76de651ef669c9af5e77b27a749b0dba3" {
		t.Fatalf("first sha should be 612120e76de651ef669c9af5e77b27a749b0dba3 but received %s", tags[1].Commit.SHA)
	}
	// Spot check a release name.
	if tags[0].Name != "v1.0.0-beta1" {
		t.Fatalf("first name should be v1.0.0-beta1 but received %s", tags[0].Name)
	}
}

func (e *testEnv) TestLatestTag(t *testing.T) {
	tag, err := e.githubClient.LatestTag("hashicorp", "vault")
	if err != nil {
		t.Fatal(err)
	}
	if tag.Name != "v1.0.0-beta1" {
		t.Fatalf("expected v1.0.0-beta1 but received %s", tag.Name)
	}
}

const tagsListResponse = `[
  {
    "name": "v1.0.0-beta1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v1.0.0-beta1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v1.0.0-beta1",
    "commit": {
      "sha": "ebc733f4ca5d362fdfb302ac75953228585c54a2",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/ebc733f4ca5d362fdfb302ac75953228585c54a2"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djEuMC4wLWJldGEx"
  },
  {
    "name": "v0.11.4",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.11.4",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.11.4",
    "commit": {
      "sha": "612120e76de651ef669c9af5e77b27a749b0dba3",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/612120e76de651ef669c9af5e77b27a749b0dba3"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTEuNA=="
  },
  {
    "name": "v0.11.3",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.11.3",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.11.3",
    "commit": {
      "sha": "fb601237bfbe4bc16ff679f642248ee8a86e627b",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/fb601237bfbe4bc16ff679f642248ee8a86e627b"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTEuMw=="
  },
  {
    "name": "v0.11.2",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.11.2",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.11.2",
    "commit": {
      "sha": "2b1a4304374712953ff606c6a925bbe90a4e85dd",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/2b1a4304374712953ff606c6a925bbe90a4e85dd"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTEuMg=="
  },
  {
    "name": "v0.11.1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.11.1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.11.1",
    "commit": {
      "sha": "8575f8fedcf8f5a6eb2b4701cb527b99574b5286",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/8575f8fedcf8f5a6eb2b4701cb527b99574b5286"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTEuMQ=="
  },
  {
    "name": "v0.11.0",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.11.0",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.11.0",
    "commit": {
      "sha": "87492f9258e0227f3717e3883c6a8be5716bf564",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/87492f9258e0227f3717e3883c6a8be5716bf564"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTEuMA=="
  },
  {
    "name": "v0.11.0-beta1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.11.0-beta1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.11.0-beta1",
    "commit": {
      "sha": "3cc78f54c664ae3833f71a95cc27d1fc9a83ad24",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/3cc78f54c664ae3833f71a95cc27d1fc9a83ad24"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTEuMC1iZXRhMQ=="
  },
  {
    "name": "v0.10.4",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.10.4",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.10.4",
    "commit": {
      "sha": "e21712a687889de1125e0a12a980420b1a4f72d3",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/e21712a687889de1125e0a12a980420b1a4f72d3"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTAuNA=="
  },
  {
    "name": "v0.10.3",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.10.3",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.10.3",
    "commit": {
      "sha": "533003e27840d9646cb4e7d23b3a113895da1dd0",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/533003e27840d9646cb4e7d23b3a113895da1dd0"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTAuMw=="
  },
  {
    "name": "v0.10.2",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.10.2",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.10.2",
    "commit": {
      "sha": "3ee0802ed08cb7f4046c2151ec4671a076b76166",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/3ee0802ed08cb7f4046c2151ec4671a076b76166"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTAuMg=="
  },
  {
    "name": "v0.10.1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.10.1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.10.1",
    "commit": {
      "sha": "756fdc4587350daf1c65b93647b2cc31a6f119cd",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/756fdc4587350daf1c65b93647b2cc31a6f119cd"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTAuMQ=="
  },
  {
    "name": "v0.10.0",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.10.0",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.10.0",
    "commit": {
      "sha": "5dd7f25f5c4b541f2da62d70075b6f82771a650d",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/5dd7f25f5c4b541f2da62d70075b6f82771a650d"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTAuMA=="
  },
  {
    "name": "v0.10.0-rc1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.10.0-rc1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.10.0-rc1",
    "commit": {
      "sha": "3890f84689e4136d73a4fb98d0795a0cea7772bd",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/3890f84689e4136d73a4fb98d0795a0cea7772bd"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuMTAuMC1yYzE="
  },
  {
    "name": "v0.9.6",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.9.6",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.9.6",
    "commit": {
      "sha": "7e1fbde40afee241f81ef08700e7987d86fc7242",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/7e1fbde40afee241f81ef08700e7987d86fc7242"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOS42"
  },
  {
    "name": "v0.9.5",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.9.5",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.9.5",
    "commit": {
      "sha": "36edb4d42380d89a897e7f633046423240b710d9",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/36edb4d42380d89a897e7f633046423240b710d9"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOS41"
  },
  {
    "name": "v0.9.4",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.9.4",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.9.4",
    "commit": {
      "sha": "2e2c89a3d8d2f4876c64dee6ba3d4a5e08691aee",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/2e2c89a3d8d2f4876c64dee6ba3d4a5e08691aee"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOS40"
  },
  {
    "name": "v0.9.3",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.9.3",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.9.3",
    "commit": {
      "sha": "5acd6a21d5a69ab49d0f7c0bf540123a9b2c696d",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/5acd6a21d5a69ab49d0f7c0bf540123a9b2c696d"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOS4z"
  },
  {
    "name": "v0.9.2",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.9.2",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.9.2",
    "commit": {
      "sha": "3bf8733cd69bb0ac14da9aaa6135bcb7f710cc5f",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/3bf8733cd69bb0ac14da9aaa6135bcb7f710cc5f"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOS4y"
  },
  {
    "name": "v0.9.1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.9.1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.9.1",
    "commit": {
      "sha": "87b6919dea55da61d7cd444b2442cabb8ede8ab1",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/87b6919dea55da61d7cd444b2442cabb8ede8ab1"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOS4x"
  },
  {
    "name": "v0.9.0",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.9.0",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.9.0",
    "commit": {
      "sha": "bdac1854478538052ba5b7ec9a9ec688d35a3335",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/bdac1854478538052ba5b7ec9a9ec688d35a3335"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOS4w"
  },
  {
    "name": "v0.8.3",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.8.3",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.8.3",
    "commit": {
      "sha": "6b29fb2b7f70ed538ee2b3c057335d706b6d4e36",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/6b29fb2b7f70ed538ee2b3c057335d706b6d4e36"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOC4z"
  },
  {
    "name": "v0.8.2",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.8.2",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.8.2",
    "commit": {
      "sha": "9afe7330e06e486ee326621624f2077d88bc9511",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/9afe7330e06e486ee326621624f2077d88bc9511"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOC4y"
  },
  {
    "name": "v0.8.1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.8.1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.8.1",
    "commit": {
      "sha": "8d76a41854608c547a233f2e6292ae5355154695",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/8d76a41854608c547a233f2e6292ae5355154695"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOC4x"
  },
  {
    "name": "v0.8.0",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.8.0",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.8.0",
    "commit": {
      "sha": "af63d879130d2ee292f09257571d371100a513eb",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/af63d879130d2ee292f09257571d371100a513eb"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOC4w"
  },
  {
    "name": "v0.8.0-rc1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.8.0-rc1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.8.0-rc1",
    "commit": {
      "sha": "eec563e5f5fda190538c14c64c8317d6ead84c13",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/eec563e5f5fda190538c14c64c8317d6ead84c13"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOC4wLXJjMQ=="
  },
  {
    "name": "v0.8.0-beta1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.8.0-beta1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.8.0-beta1",
    "commit": {
      "sha": "3fb82dbb6620f1ea798bb67bfe224e28de19f754",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/3fb82dbb6620f1ea798bb67bfe224e28de19f754"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuOC4wLWJldGEx"
  },
  {
    "name": "v0.7.3",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.7.3",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.7.3",
    "commit": {
      "sha": "0b20ae0b9b7a748d607082b1add3663a28e31b68",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/0b20ae0b9b7a748d607082b1add3663a28e31b68"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuNy4z"
  },
  {
    "name": "v0.7.2",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.7.2",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.7.2",
    "commit": {
      "sha": "d28dd5a018294562dbc9a18c95554d52b5d12390",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/d28dd5a018294562dbc9a18c95554d52b5d12390"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuNy4y"
  },
  {
    "name": "v0.7.1",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.7.1",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.7.1",
    "commit": {
      "sha": "1767ced3a336aee9dbce7f4f0ad652ff44e8ddab",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/1767ced3a336aee9dbce7f4f0ad652ff44e8ddab"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuNy4x"
  },
  {
    "name": "v0.7.0",
    "zipball_url": "https://api.github.com/repos/hashicorp/vault/zipball/v0.7.0",
    "tarball_url": "https://api.github.com/repos/hashicorp/vault/tarball/v0.7.0",
    "commit": {
      "sha": "614deacfca3f3b7162bbf30a36d6fc7362cd47f0",
      "url": "https://api.github.com/repos/hashicorp/vault/commits/614deacfca3f3b7162bbf30a36d6fc7362cd47f0"
    },
    "node_id": "MDM6UmVmMzEyODg5NTg6djAuNy4w"
  }
]`
