package pkg

type TokenResponse struct {
	RESP struct {
		NOTICES        []string
		WARNINGS       []string
		ERROR          []string
		API_ACCESS_KEY string
	} `json:"resp"`
}

type SignedUrlResponse struct {
	Resp struct {
		PRESIGNED_URL string
	} `json:"resp"`
}

type JobReleaseResponse struct {
	RELEASE_NAME string
	Latest       int
	Active       int
}

type ServiceReleaseResponse struct {
	RELEASE_NAME string
}

type ReleasesResponse struct {
	Resp struct {
		Notices  []string
		Warnings []string
		Errors   []string
		Jobs     []JobReleaseResponse
		Services []ServiceReleaseResponse
	}
}

type LogResponse struct {
	Resp struct {
		Notices  []string
		Warnings []string
		Errors   []string
		Logs     []string
	}
}

type JobStatus struct {
	Resp struct {
		Notices    []string
		Warnings   []string
		Errors     []string
		JOB_STATUS struct {
			Reason     string
			ExitCode   int
			StartedAt  string
			FinishedAt string
		}
	}
}

type ServiceStatus struct {
	Resp struct {
		Notices  []string
		Warnings []string
		Errors   []string
		Response struct {
			CreationTimestamp string
			Status            string
		}
	}
}

type DownResponse struct {
	Resp struct {
		Notices  []string
		Warnings []string
		Errors   []string
	}
}

type UpRelease struct {
	Notices          []string
	Warnings         []string
	Errors           []string
	COMMAND_RESPONSE struct {
		FIRST_DEPLOYED string
		LAST_DELOYED   string
		Deleted        string
		Description    string
		Status         string
		Notes          []string
		Action         string
		RELEASE_NAME   string
	}
}

type UpResponse struct {
	Resp map[string]UpRelease
}

type Response struct {
	Resp struct {
		Msg      string
		Complete string
		ASK_FOR  string
		Options  []string
		Examples []string
		Config   struct {
			PROJECT_NAME string
			Service      map[string]Service
		}
	} `json:"resp"`
}

type Image struct {
	Name string
	Tag  string
}

type Service struct {
	Image       Image
	Command     string
	Ports       []string
	Environment []struct {
		Key string
		Val string
	}
	WorkingDir string
	VALUES     map[string]interface{}
}

type Project struct {
	PROJECT_NAME string
	SERVICES     []map[string]Service
}
