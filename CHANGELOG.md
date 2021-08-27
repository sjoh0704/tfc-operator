# tfc-operator changelog!!
All notable changes to this project will be documented in this file.

<!-------------------- v5.0.25.0 start -------------------->

## tfc-operator_5.0.25.0 (2021. 08. 27. (금) 18:58:39 KST)

### Added

### Changed
  - [mod] Private Git 인증처리 로직 수정 (ID/PW -> Personal Access Token) by gyeongyeol-choi

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.25.0 end --------------------->

<!-------------------- v5.0.24.0 start -------------------->

## tfc-operator_5.0.24.0 (2021. 08. 19. (목) 17:04:34 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.24.0 end --------------------->

<!-------------------- v5.0.23.0 start -------------------->

## tfc-operator_5.0.23.0 (2021. 08. 12. (목) 13:07:38 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.23.0 end --------------------->

<!-------------------- v5.0.22.0 start -------------------->

## tfc-operator_5.0.22.0 (2021. 08. 05. (목) 15:18:09 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.22.0 end --------------------->

<!-------------------- v5.0.21.0 start -------------------->

## tfc-operator_5.0.21.0 (2021. 07. 29. (목) 15:36:02 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc
  - [etc] 코드 리팩토링 (전역변수 처리, 공통로직 함수 모듈화 등) by gyeongyeol-choi

<!--------------------- v5.0.21.0 end --------------------->

<!-------------------- v5.0.20.0 start -------------------->

## tfc-operator_5.0.20.0 (2021. 07. 22. (목) 14:21:22 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.20.0 end --------------------->

<!-------------------- v5.0.19.0 start -------------------->

## tfc-operator_5.0.19.0 (2021. 07. 15. (목) 17:23:08 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.19.0 end --------------------->

<!-------------------- v5.0.18.0 start -------------------->

## tfc-operator_5.0.18.0 (2021. 07. 09. (금) 09:40:27 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.18.0 end --------------------->

<!-------------------- v5.0.17.7 start -------------------->

## tfc-operator_5.0.17.7 (2021. 07. 08. (목) 18:09:31 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc
  - [etc] Terraform Worker Pod 생성 지연 시간 변경 (5 sec -> 15 sec) by gyeongyeol-choi
  - [etc] CRD 업데이트에 따른 key mapping 파일 수정 by gyeongyeol-choi

<!--------------------- v5.0.17.7 end --------------------->

<!-------------------- v5.0.17.6 start -------------------->

## tfc-operator_5.0.17.6 (2021. 07. 08. (목) 13:33:28 KST)

### Added

### Changed
  - [mod] Terraform Plan/Apply/Destroy 명령 처리 과정에서 에러 발생 시 로그 내용 조회 가능하도록 로직 수정 by gyeongyeol-choi

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.17.6 end --------------------->

<!-------------------- v5.0.17.4 start -------------------->

## tfc-operator_5.0.17.4 (2021. 07. 08. (목) 12:36:38 KST)

### Added

### Changed
  - [mod] Terraform Version에 따른 Plugin 인증 처리 로직 추가 by gyeongyeol-choi

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.17.4 end --------------------->

<!-------------------- v5.0.17.3 start -------------------->

## tfc-operator_5.0.17.3 (2021. 07. 05. (월) 19:27:40 KST)

### Added

### Changed
  - [mod] Terraform Init 중복 처리 시 에러 수정 by gyeongyeol-choi

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.17.3 end --------------------->

<!-------------------- v5.0.17.2 start -------------------->

## tfc-operator_5.0.17.2 (2021. 07. 05. (월) 17:56:39 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc
  - [etc] Default Branch (main, master)로의  Checkout 실패에 대한 예외 처리 (already exists) by gyeongyeol-choi

<!--------------------- v5.0.17.2 end --------------------->

<!-------------------- v5.0.17.1 start -------------------->

## tfc-operator_5.0.17.1 (2021. 07. 05. (월) 16:31:50 KST)

### Added

### Changed

### Fixed
  - [ims][265432] Git Credential ID/PW 내 특수문자에 대한 예외처리 by gyeongyeol-choi

### CRD yaml

### Etc

<!--------------------- v5.0.17.1 end --------------------->

<!-------------------- v5.0.17.0 start -------------------->

## tfc-operator_5.0.17.0 (2021. 07. 01. (목) 17:48:16 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.17.0 end --------------------->

<!-------------------- v5.0.16.3 start -------------------->

## tfc-operator_5.0.16.3 (2021. 07. 01. (목) 15:14:17 KST)

### Added

### Changed
  - [mod] 에러 발생에 대한 이전상태 정보 관리 기능 추가 by gyeongyeol-choi
  - [mod] repo type 혼동에 따른  무한 프리징 에러 해결 by gyeongyeol-choi

### Fixed

### CRD yaml

### Etc
  - [etc] 에러 상세 내용을 명확히 파악 가능하도록 Reason 필드 재정의 by gyeongyeol-choi

<!--------------------- v5.0.16.3 end --------------------->

<!-------------------- v5.0.16.2 start -------------------->

## tfc-operator_5.0.16.2 (2021. 06. 25. (금) 18:07:33 KST)

### Added

### Changed
  - [mod] Secret 참조에 대한 에러 처리 로직 세분화 by gyeongyeol-choi

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.16.2 end --------------------->

<!-------------------- v5.0.16.1 start -------------------->

## tfc-operator_5.0.16.1 (2021. 06. 25. (금) 17:00:29 KST)

### Added

### Changed

### Fixed
  - [ims][262346] private repo type의 Claim 생성 시, secret (git credential)을 넣지 않았을 경우의 에러 처리 로직 구현 by gyeongyeol-choi

### CRD yaml

### Etc

<!--------------------- v5.0.16.1 end --------------------->

<!-------------------- v5.0.16.0 start -------------------->

## tfc-operator_5.0.16.0 (2021. 06. 25. (금) 11:15:16 KST)

### Added

### Changed
  - [mod] terraform apply 명령 처리 시, Git 관련 정보 저장 by gyeongyeol-choi
  - [mod] Public Repo 지원을 위한 Git Credential Secret 처리 로직 수정 by gyeongyeol-choi

### Fixed

### CRD yaml

### Etc

<!--------------------- v5.0.16.0 end --------------------->

<!-------------------- v5.0.15.0 start -------------------->

## tfc-operator_5.0.15.0 (2021. 06. 17. (목) 15:14:33 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc
  - [etc] TFApplyClaim Spec 중 일부 Required Field 처리 by gyeongyeol-choi

<!--------------------- v5.0.15.0 end --------------------->

<!-------------------- v5.0.14.2 start -------------------->

## tfc-operator_5.0.14.2 (2021. 06. 11. (금) 17:17:09 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc
  - [etc] Changlog 스크립트 수정 by gyeongyeol-choi

<!--------------------- v5.0.14.2 end --------------------->

<!-------------------- v5.0.14.1 start -------------------->

## tfc-operator_5.0.14.1 (2021. 06. 11. (금) 16:59:09 KST)

### Added

### Changed
  - [mod] TFApplyClaim 승인 과정 Reject 처리 로직 추가 by gyeongyeol-choi
  - [mod] terraform init 시 플러그인 모듈 인증 skip 로직 추가 by gyeongyeol-choi

### Fixed
  - [ims][262548] tfc-worker 이미지 경로 변경 가능하도록 Deployment Environment 로직 추가 by gyeongyeol-choi

### CRD yaml

### Etc
  - [etc] Manifest 경로 수정 by gyeongyeol-choi
  - [etc] key-mapping 파일 추가 by gyeongyeol-choi
  - [etc] update image name & tag in kustomization by gyeongyeol-choi
  - [etc] init by gyeongyeol-choi

<!--------------------- v5.0.14.1 end --------------------->

<!-------------------- v5.0.13.2 start -------------------->

## tfc-operator_5.0.13.2 (2021. 06. 09. (수) 16:58:24 KST)

### Added

### Changed
  - [mod] terraform init 시 플러그인 모듈 인증 skip 로직 추가 by gyeongyeol-choi

### Fixed
  - [ims][262548] tfc-worker 이미지 경로 변경 가능하도록 Deployment Environment 로직 추가 by gyeongyeol-choi

### CRD yaml

### Etc
  - [etc] Manifest 경로 수정 by gyeongyeol-choi
  - [etc] key-mapping 파일 추가 by gyeongyeol-choi
  - [etc] update image name & tag in kustomization by gyeongyeol-choi
  - [etc] init by gyeongyeol-choi

<!--------------------- v5.0.13.2 end --------------------->

<!-------------------- v5.0.12.1 start -------------------->

## tfc-operator_5.0.12.1 (2021. 06. 02. (수) 14:05:50 KST)

### Added

### Changed

### Fixed
  - [ims][262548] tfc-worker 이미지 경로 변경 가능하도록 Deployment Environment 로직 추가 by gyeongyeol-choi

### CRD yaml

### Etc
  - [etc] Manifest 경로 수정 by gyeongyeol-choi
  - [etc] key-mapping 파일 추가 by gyeongyeol-choi
  - [etc] update image name & tag in kustomization by gyeongyeol-choi
  - [etc] init by gyeongyeol-choi

<!--------------------- v5.0.12.1 end --------------------->

<!-------------------- v5.0.1.1 start -------------------->

## tfc-operator_5.0.1.1 (2021. 05. 27. (목) 18:04:35 KST)

### Added

### Changed

### Fixed

### CRD yaml

### Etc
  - [etc] key-mapping 파일 추가 by gyeongyeol-choi
  - [etc] update image name & tag in kustomization by gyeongyeol-choi
  - [etc] init by gyeongyeol-choi

<!--------------------- v5.0.1.1 end --------------------->

<!-------------------- v5.0.1.0 start -------------------->

## tfc-operator_5.0.1.0 (Thu May 20 03:32:31 KST 2021)

### Added

### Changed

### Fixed

### CRD yaml

### Etc
  - [etc] update image name & tag in kustomization by gyeongyeol-choi
  - [etc] init by gyeongyeol-choi

<!--------------------- v5.0.1.0 end --------------------->
