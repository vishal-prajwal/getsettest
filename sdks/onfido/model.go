package onfido

type Location struct {
	IPAddress          string `json:"ip_address"`
	CountryOfResidence string `json:"country_of_residence"`
}

type CreateApplicantRequest struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Location  Location `json:"location"`
}

type CreateApplicantResponse struct {
	ID    string      `json:"id"`
	Error interface{} `json:"error"`
}

type CreateCheckRequest struct {
	ApplicantId string   `json:"applicant_id"`
	ReportNames []string `json:"report_names"`
}

type CreateCheckResponse struct {
	ID        string      `json:"id"`
	ReportIds []string    `json:"report_ids"`
	Status    string      `json:"status"`
	Error     interface{} `json:"error"`
}

type UploadDocumentResponse struct {
	ID           string      `json:"id"`
	FileName     string      `json:"file_name"`
	FileType     string      `json:"file_type"`
	Type         string      `json:"type"`
	Size         int         `json:"file_size"`
	Side         string      `json:"side"`
	DownloadHref string      `json:"download_href"`
	Error        interface{} `json:"error"`
}

type ReportResponse struct {
	CheckId    string `json:"check_id"`
	CreatedAt  string `json:"created_at"`
	ReportId   string `json:"id"`
	Name       string `json:"name"`
	Properties struct {
		DOB             string `json:"date_of_birth"`
		DateOfExpiry    string `json:"date_of_expiry"`
		DocumentNumbers []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"document_numbers"`
		DocType        string `json:"document_type"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Gender         string `json:"gender"`
		IssuingCountry string `json:"issuing_country"`
		Nationality    string `json:"nationality"`
		Address        string `json:"address"`
		AddressLines   struct {
			City          string `json:"city"`
			Country       string `json:"country"`
			PostalCode    string `json:"postal_code"`
			State         string `json:"state"`
			StreetAddress string `json:"street_address"`
		} `json:"address_lines"`
	} `json:"properties"`
	Result    string `json:"result"`
	Breakdown struct {
		DataComparison struct {
			Result    interface{} `json:"result"`
			Breakdown struct {
				DateOfExpiry struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"date_of_expiry"`
				IssuingCountry struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"issuing_country"`
				DocumentType struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"document_type"`
				DocumentNumbers struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"document_numbers"`
				Gender struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"gender"`
				DateOfBirth struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"date_of_birth"`
				LastName struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"last_name"`
				FirstName struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"first_name"`
			} `json:"breakdown"`
		} `json:"data_comparison"`
		DataConsistency struct {
			Result    interface{} `json:"result"`
			Breakdown struct {
				MultipleDataSourcesPresent struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"multiple_data_sources_present"`
				DateOfExpiry struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"date_of_expiry"`
				DocumentType struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"document_type"`
				Nationality struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"nationality"`
				IssuingCountry struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"issuing_country"`
				DocumentNumbers struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"document_numbers"`
				Gender struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"gender"`
				DateOfBirth struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"date_of_birth"`
				LastName struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"last_name"`
				FirstName struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"first_name"`
			} `json:"breakdown"`
		} `json:"data_consistency"`
		CompromisedDocument struct {
			Result string `json:"result"`
		} `json:"compromised_document"`
		AgeValidation struct {
			Result    string `json:"result"`
			Breakdown struct {
				MinimumAcceptedAge struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"minimum_accepted_age"`
			} `json:"breakdown"`
		} `json:"age_validation"`
		ImageIntegrity struct {
			Result    string `json:"result"`
			Breakdown struct {
				SupportedDocument struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"supported_document"`
				ImageQuality struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"image_quality"`
				ColourPicture struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"colour_picture"`
				ConclusiveDocumentQuality struct {
					Result     string `json:"result"`
					Properties struct {
						MissingBack                  string `json:"missing_back"`
						DigitalDocument              string `json:"digital_document"`
						PuncturedDocument            string `json:"punctured_document"`
						CornerRemoved                string `json:"corner_removed"`
						WatermarksDigitalTextOverlay string `json:"watermarks_digital_text_overlay"`
						AbnormalDocumentFeatures     string `json:"abnormal_document_features"`
						ObscuredSecurityFeatures     string `json:"obscured_security_features"`
						ObscuredDataPoints           string `json:"obscured_data_points"`
					} `json:"properties"`
				} `json:"conclusive_document_quality"`
			} `json:"breakdown"`
		} `json:"image_integrity"`
		PoliceRecord struct {
			Result interface{} `json:"result"`
		} `json:"police_record"`
		IssuingAuthority struct {
			Result    interface{} `json:"result"`
			Breakdown struct {
				NfcPassiveAuthentication struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"nfc_passive_authentication"`
				NfcActiveAuthentication struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"nfc_active_authentication"`
			} `json:"breakdown"`
		} `json:"issuing_authority"`
		DataValidation struct {
			Result    string `json:"result"`
			Breakdown struct {
				DocumentNumbers struct {
					Result     string `json:"result"`
					Properties struct {
						DocumentNumber string `json:"document_number"`
					} `json:"properties"`
				} `json:"document_numbers"`
				Gender struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"gender"`
				DateOfBirth struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"date_of_birth"`
				DocumentExpiration struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"document_expiration"`
				ExpiryDate struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"expiry_date"`
				Mrz struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"mrz"`
				Barcode struct {
					Result     interface{} `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"barcode"`
			} `json:"breakdown"`
		} `json:"data_validation"`
		VisualAuthenticity struct {
			Result    string `json:"result"`
			Breakdown struct {
				Other struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"other"`
				PictureFaceIntegrity struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"picture_face_integrity"`
				FaceDetection struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"face_detection"`
				DigitalTampering struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"digital_tampering"`
				OriginalDocumentPresent struct {
					Result     string `json:"result"`
					Properties struct {
						Scan                   string `json:"scan"`
						DocumentOnPrintedPaper string `json:"document_on_printed_paper"`
						Screenshot             string `json:"screenshot"`
						PhotoOfScreen          string `json:"photo_of_screen"`
					} `json:"properties"`
				} `json:"original_document_present"`
				SecurityFeatures struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"security_features"`
				Fonts struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"fonts"`
				Template struct {
					Result     string `json:"result"`
					Properties struct {
					} `json:"properties"`
				} `json:"template"`
			} `json:"breakdown"`
		} `json:"visual_authenticity"`
	}
	Error interface{} `json:"error"`
}
