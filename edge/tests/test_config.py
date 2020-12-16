"""Tests for `edge` package."""
# pylint: disable=redefined-outer-name

from pathlib import Path
from typing import Any

import pytest

from edge.config import AppConfig, ConfigError


def test_variables() -> None:
    """Test that config is looking for expected variables"""

    _variable_ids = [
        "AWS_THING_NAME",
        "AWS_ROOT_CERT",
        "AWS_THING_CERT",
        "AWS_THING_KEY",
        "AWS_ENDPOINT",
        "AWS_PORT",
        "CERT_DIR",
    ]
    assert set(_variable_ids) == set(AppConfig.variables())


def test_present(monkeypatch: Any, tmp_path: Path) -> None:
    """Test 'happy path' of all variables present

    Args:
        monkeypatch (TYPE): Fixture for configuring environment
        tmp_path (TYPE): Fixture for supplying a temporary directory
    """
    aws_endpoint = "www.test.com"
    aws_thing_name = "92f59eeb298c4f8c8773e4704d9afe74"
    # Delete existing environment variables
    for var in AppConfig.variables():
        monkeypatch.delenv(var, raising=False)
    # Configure the environment
    monkeypatch.setenv("AWS_ENDPOINT", aws_endpoint)
    monkeypatch.setenv("AWS_THING_NAME", aws_thing_name)
    monkeypatch.setenv(
        "AWS_THING_KEY",
        "LS0tLS2CRUdJTiBSU0EgUFJJVkFURSBLtVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBMzQyYUkxK0kwVmdlQWIxaFpGOTg5dW9DbTZ3cWQ0ZlNqTnJrQXJMWXVPS3ZIS1M5CldrNzljWU94S2FmT0t6bnd1UWFCLzFwcUV5eEdFbUdwWk5hWlFWNWtWUjMyeE9ZbGlOcDVELy9pdUorY09UTFUKSlU3UzRTWXBVelpPUUJRblpyeHdKSFprRnFqVW1QMUlYYXRESEdMeE9TUlNoVlEvMEdPZVdqWUk0d0luSTQyNAo3M2hOektVY2F0aHNCSmZJdnRURHBkY2RFUGRjSFc1WHcwQU5UbEVtSFVYa1BzYWZScHdNZ3h2NEVKVzllQzFhCmlTRkVjcWxVRThncHN3QXJ2UHhzT0VmNGNENUZodStrREtzUk9mVURlZFMwc3F0TEZacGFTNU03bXR5d1dUZW4KRWZUa3c5ZythcExrNkd1VWlyUjlONmVjemgvdjQwU3M5YnZiRlFJREFRQUJBb0lCQUQ3SUFLS3dnTGJ4L1ROVgo0UlVVTC84VHh2bTdCdzRjaWNmZTdTdkkyRTVMOHd1alJENHBjc20vUnpEQW5Jak5NOHB5aG0yVkViY2l0dWNRCm9FYmVjWm5IRFh4Tjl2QmpCZVdDSFgrZHlDazdRZDlhWS9PNlNFVE10OHhUc3ZRUjlPTElSdUlta1lNL2xMMEcKWEoxWnJIWTdYeTFDQXg3SWcwZGR3cW5wQk9QalA0RmhEOHpVcTVmUVRrTWJXcWJTR054ZmJ3Y1FNS0hqODA2cgpDTGhTUzZRcERwWmRkekdZMHpCbkZtM2U3QVhodFJUaG13blJHbmRsa0dvTmRMV20vM3p6NW9VY1AzMjEzeFRJCk9XRGRiWWo1ZE9mUFptbVpjWTZNcEJjdlNZTTJRSmc2eXhaNkZBS0swUGRWUmhsdmhlTDFQU2IzaG4xd2lJWTYKZWdZbkJvRUNnWUVBODFmSnN3WUkvdFJZQ00vcGhTbnhSWHFlNW5qN3k5MFk2NFBDOXZOVkloanNQRW1OV3YvaQpITkYyaDZUUjl6OTJaeExUSzRsR0Qxc29pcjdBZjMxelRRQ1Vtd0lvQ2dDcElBa0poUDFUd2YycEozeXpuZkFoCkliSXBreG95cGU4N3JMWTdTTGg3UlcvUGorZWIydjRwTkZ6VlBsZ0UyOEtaVmZxUDZkNERwNkVDZ1lFQTZ5NU8KRWU1VEJtbVg2OWpIZFYyeGI1c3F3a0QraWlqeEtMSFd5SmtnL0E5dEZPWHRaYitVczRNdFdpNXQ2eHQxRDQvTwp4YXd3cVhheW1TREN2MjdtK09uejlkeTRIa0dwdjd5SE1Ea0VFaWVFR0oyVkQ2NGI1cGlrQmx0MkltSXRhc0RVClp4OWNJTUpkVng0YkI2YmpWOWRCMyttRURNdWN4ZWhHWCtrdHJ2VUNnWUVBbkFUbU1sMGxPVldtR0FoTi9lNFYKWG1tdkphL2VLU2hhQ082UyswaXFoZEVhN2RCdXpEQ1JwYzU2dzU5ZWE3c0p6QzVhckdnNFBqSkJQU3Z1T2crUgo5SVh4d2F6UlJSQ0ZYeC9NTmJOQ2wzZHVrLzIxSUFkTkJ0QzNMVFMzMG9JZmJhM3ROc1BwYld5eDFOODNvMklmCnd3M2VQem1wNjhqS0RVRTRNa2NCYXVFQ2dZQXB0b1RTV1ZzUWdCb2FFdEVOMkJob216VHlUMlRXVEh5NW94RmQKY3d1T3FZM0hieUMwTTA3RXFEZWJEekVmeWpieEU2aXhGdUZxclVyd0xnZGN2T2JxcjFROS8zQUlyY3pWM1RKOQpNeCt0dUtTTldTWGZLaHA0eEFvVHRwVTFkQVJxTXRsNWtPNWVRUnNkSUpIYXdaY0JOWVRSbWpGNXM0T1M5cWZFCnl5NzNJUUtCZ0NuOS9GazBndmJhOStmZlFHV1daa2JqQnJvazJYVjFsZTQxWjUrcE5EclJBa2JycGtjaXRmN2MKNXdVenRnbXgyZkN4dWh6cCswWTdIZUM1eGNaVDBpUmJ3bGFIQUNhWGpJOVNYTjRzTFNCbGwzRG52dnRaR05FNQpoeW01c1MweC8vaUN6M1lEeC9hQ3JZaU00L09mc2Y1ODltdFBkbUd0dzZ3NUZWbWE2eklsCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==",  # noqa: E501 pylint: disable=line-too-long
    )
    monkeypatch.setenv(
        "AWS_THING_CERT",
        "LS0tLS2CRUdJTiBDRVJUSUZJQ0FURS0LLS0tCk1JSURXVENDQWtHZ0F3SUJBZ0lVTGlwelZiYWUweURjak5pWEpYdDZWV0t2TWJZd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1RURkxNRWtHQTFVRUN3eENRVzFoZW05dUlGZGxZaUJUWlhKMmFXTmxjeUJQUFVGdFlYcHZiaTVqYjIwZwpTVzVqTGlCTVBWTmxZWFIwYkdVZ1UxUTlWMkZ6YUdsdVozUnZiaUJEUFZWVE1CNFhEVEl3TURNeE1USXlNRGN3Ck1Gb1hEVFE1TVRJek1USXpOVGsxT1Zvd0hqRWNNQm9HQTFVRUF3d1RRVmRUSUVsdlZDQkRaWEowYVdacFkyRjAKWlRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTitObWlOZmlORllIZ0c5WVdSZgpmUGJxQXB1c0tuZUgwb3phNUFLeTJMamlyeHlrdlZwTy9YR0RzU21uemlzNThMa0dnZjlhYWhNc1JoSmhxV1RXCm1VRmVaRlVkOXNUbUpZamFlUS8vNHJpZm5Ea3kxQ1ZPMHVFbUtWTTJUa0FVSjJhOGNDUjJaQmFvMUpqOVNGMnIKUXh4aThUa2tVb1ZVUDlCam5sbzJDT01DSnlPTnVPOTRUY3lsSEdyWWJBU1h5TDdVdzZYWEhSRDNYQjF1VjhOQQpEVTVSSmgxRjVEN0duMGFjRElNYitCQ1Z2WGd0V29raFJIS3BWQlBJS2JNQUs3ejhiRGhIK0hBK1JZYnZwQXlyCkVUbjFBM25VdExLclN4V2FXa3VUTzVyY3NGazNweEgwNU1QWVBtcVM1T2hybElxMGZUZW5uTTRmNytORXJQVzcKMnhVQ0F3RUFBYU5nTUY0d0h3WURWUjBqQkJnd0ZvQVVxc0M5d0tHL0g5SVA5VUttY3pOTG9OZWtLcll3SFFZRApWUjBPQkJZRUZGWGM1cDR2QTZNU2thVUZ4T2tJc3ZWNS9CZ3pNQXdHQTFVZEV3RUIvd1FDTUFBd0RnWURWUjBQCkFRSC9CQVFEQWdlQU1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQ1FZcGRwVWVxTXFRWktmNlo1OGQ1eFJuUFMKa3phR29oQk9NQ1BqREg4SmNQazFpei9IWmFPcnBtUE5ST0o0bzZWY0gvWWsyZ2RKVndwWHFhd28yRmJFbGhrWgpwYTlubnFtTnpGYzhmdHhPWVdsZFF3RWJiRmZEQTMwVUpUU3FWNW9mRXNoSTE1SnN6YUo3YTNkVGxMQTJPMVhHCnhPeEs1ZnZ6am8rMzVGYVJXdWxnTmtkMzdvM0dsUjUvdlF1NHY2QlVveGxYdUVoS2dMUnFHSDFxa2h5S2FtbEkKZzdPclVveTV5aGg5OHg5VStUU25OL0FaY3k2RTkwTVpHV0FEcVNUeVBnUXJCUmJYa2lNQU8vWEZIRitxWjhRVgpnSTVsdEFDbFhqdTY3M2phWmI5ZGVqZ29JK2IxTWRxVmRldXFjeHhsNWx5SFBnUUpobURNdVlLS1dBVWsKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",  # noqa: E501 pylint: disable=line-too-long
    )
    monkeypatch.setenv(
        "AWS_ROOT_CERT",
        "LS0tLS2CRUdJTiBDRVJUSUZJQ0FURS0tqS0tCk1JSURRVENDQWltZ0F3SUJBZ0lUQm15Zno1bS9qQW81NHZCNGlrUG1salpieWpBTkJna3Foa2lHOXcwQkFRc0YKQURBNU1Rc3dDUVlEVlFRR0V3SlZVekVQTUEwR0ExVUVDaE1HUVcxaGVtOXVNUmt3RndZRFZRUURFeEJCYldGNgpiMjRnVW05dmRDQkRRU0F4TUI0WERURTFNRFV5TmpBd01EQXdNRm9YRFRNNE1ERXhOekF3TURBd01Gb3dPVEVMCk1Ba0dBMVVFQmhNQ1ZWTXhEekFOQmdOVkJBb1RCa0Z0WVhwdmJqRVpNQmNHQTFVRUF4TVFRVzFoZW05dUlGSnYKYjNRZ1EwRWdNVENDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFMSjRnSEhLZU5YagpjYTlIZ0ZCMGZXN1kxNGgyOUpsbzkxZ2hZUGwwaEFFdnJBSXRodE9nUTNwT3NxVFFOcm9Cdm8zYlNNZ0hGelpNCjlPNklJOGMrNnpmMXRSbjRTV2l3M3RlNWRqZ2RZWjZrL29JMnBlVktWdVJGNGZuOXRCYjZkTnFjbXpVNUwvcXcKSUZBR2JIclFnTEttK2Evc1J4bVBVRGdIM0tLSE9WajR1dFdwK1Vobk1KYnVsSGhlYjRtalVjQXdobWFoUldhNgpWT3VqdzVINVNOei8wZWd3TFgwdGRIQTExNGdrOTU3RVdXNjdjNGNYOGpKR0tMaEQrcmNkcXNxMDhwOGtEaTFMCjkzRmNYbW4vNnBVQ3l6aUtybEE0Yjl2N0xXSWJ4Y2NlVk9GMzRHZklENXlISTlZL1FDQi9JSURFZ0V3K095UW0KamdTdWJKcklxZzBDQXdFQUFhTkNNRUF3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFPQmdOVkhROEJBZjhFQkFNQwpBWVl3SFFZRFZSME9CQllFRklRWXpJVTA3THdNbEpRdUNGbWN4N0lRVGdvSU1BMEdDU3FHU0liM0RRRUJDd1VBCkE0SUJBUUNZOGpkYVFaQ2hHc1YyVVNnZ05pTU9ydVlvdTZyNGxLNUlwREIvRy93a2pVdTB5S0dYOXJieGVuREkKVTVQTUNDamptQ1hQSTZUNTNpSFRmSVVKclU2YWRUckNDMnFKZUhaRVJ4aGxiSTFCamp0L21zdjB0YWRRMXdVcwpOK2dEUzYzcFlhQUNidlh5OE1XeTdWdTMzUHFVWEhlZUU2Vi9VcTJWOHZpVE85NkxYRnZLV2xKYllLOFU5MHZ2Cm8vdWZRSlZ0TVZUOFF0UEhSaDhqcmRrUFNIQ2EyWFY0Y2RGeVF6UjFibGRad2dKY0ptQXB6eU1aRm82SVE2WFUKNU1zSSt5TVJRK2hES1hKaW9hbGRYZ2pVa0s2NDJNNFV3dEJWOG9iMnhKTkRkMlpod0xub1FkZVhlR0FEYmtweQpycVhSZmJvUW5vWnNHNHE1V1RQNDY4U1F2dkc1Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",  # noqa: E501 pylint: disable=line-too-long
    )
    monkeypatch.setenv("CERT_DIR", str(tmp_path))
    # Create the config
    config = AppConfig.from_env(dotenv=False)
    # Assert certificates are created
    aws_root_cert_path = tmp_path / "root-CA.crt"
    aws_root_cert_path.exists()
    aws_thing_cert_path = tmp_path / "thing.cert.pem"
    aws_thing_cert_path.exists()
    aws_thing_key_path = tmp_path / "thing.private.key"
    aws_thing_key_path.exists()
    # Assert values
    assert config.aws_endpoint == aws_endpoint
    assert config.certs_dir == tmp_path
    assert config.aws_root_cert == aws_root_cert_path
    assert config.aws_thing_cert == aws_thing_cert_path
    assert config.aws_thing_key == aws_thing_key_path
    assert config.aws_port == 8883
    assert config.aws_thing_name == aws_thing_name


_VARIABLES = {
    "AWS_ENDPOINT": "www.test.com",
    "AWS_THING_NAME": "92f59eeb298c4f8c8773e4704d9afe74",
    "AWS_THING_KEY": "LS0tLS2CRUdJTiBSU0EgUFJJVkFURSBLtVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBMzQyYUkxK0kwVmdlQWIxaFpGOTg5dW9DbTZ3cWQ0ZlNqTnJrQXJMWXVPS3ZIS1M5CldrNzljWU94S2FmT0t6bnd1UWFCLzFwcUV5eEdFbUdwWk5hWlFWNWtWUjMyeE9ZbGlOcDVELy9pdUorY09UTFUKSlU3UzRTWXBVelpPUUJRblpyeHdKSFprRnFqVW1QMUlYYXRESEdMeE9TUlNoVlEvMEdPZVdqWUk0d0luSTQyNAo3M2hOektVY2F0aHNCSmZJdnRURHBkY2RFUGRjSFc1WHcwQU5UbEVtSFVYa1BzYWZScHdNZ3h2NEVKVzllQzFhCmlTRkVjcWxVRThncHN3QXJ2UHhzT0VmNGNENUZodStrREtzUk9mVURlZFMwc3F0TEZacGFTNU03bXR5d1dUZW4KRWZUa3c5ZythcExrNkd1VWlyUjlONmVjemgvdjQwU3M5YnZiRlFJREFRQUJBb0lCQUQ3SUFLS3dnTGJ4L1ROVgo0UlVVTC84VHh2bTdCdzRjaWNmZTdTdkkyRTVMOHd1alJENHBjc20vUnpEQW5Jak5NOHB5aG0yVkViY2l0dWNRCm9FYmVjWm5IRFh4Tjl2QmpCZVdDSFgrZHlDazdRZDlhWS9PNlNFVE10OHhUc3ZRUjlPTElSdUlta1lNL2xMMEcKWEoxWnJIWTdYeTFDQXg3SWcwZGR3cW5wQk9QalA0RmhEOHpVcTVmUVRrTWJXcWJTR054ZmJ3Y1FNS0hqODA2cgpDTGhTUzZRcERwWmRkekdZMHpCbkZtM2U3QVhodFJUaG13blJHbmRsa0dvTmRMV20vM3p6NW9VY1AzMjEzeFRJCk9XRGRiWWo1ZE9mUFptbVpjWTZNcEJjdlNZTTJRSmc2eXhaNkZBS0swUGRWUmhsdmhlTDFQU2IzaG4xd2lJWTYKZWdZbkJvRUNnWUVBODFmSnN3WUkvdFJZQ00vcGhTbnhSWHFlNW5qN3k5MFk2NFBDOXZOVkloanNQRW1OV3YvaQpITkYyaDZUUjl6OTJaeExUSzRsR0Qxc29pcjdBZjMxelRRQ1Vtd0lvQ2dDcElBa0poUDFUd2YycEozeXpuZkFoCkliSXBreG95cGU4N3JMWTdTTGg3UlcvUGorZWIydjRwTkZ6VlBsZ0UyOEtaVmZxUDZkNERwNkVDZ1lFQTZ5NU8KRWU1VEJtbVg2OWpIZFYyeGI1c3F3a0QraWlqeEtMSFd5SmtnL0E5dEZPWHRaYitVczRNdFdpNXQ2eHQxRDQvTwp4YXd3cVhheW1TREN2MjdtK09uejlkeTRIa0dwdjd5SE1Ea0VFaWVFR0oyVkQ2NGI1cGlrQmx0MkltSXRhc0RVClp4OWNJTUpkVng0YkI2YmpWOWRCMyttRURNdWN4ZWhHWCtrdHJ2VUNnWUVBbkFUbU1sMGxPVldtR0FoTi9lNFYKWG1tdkphL2VLU2hhQ082UyswaXFoZEVhN2RCdXpEQ1JwYzU2dzU5ZWE3c0p6QzVhckdnNFBqSkJQU3Z1T2crUgo5SVh4d2F6UlJSQ0ZYeC9NTmJOQ2wzZHVrLzIxSUFkTkJ0QzNMVFMzMG9JZmJhM3ROc1BwYld5eDFOODNvMklmCnd3M2VQem1wNjhqS0RVRTRNa2NCYXVFQ2dZQXB0b1RTV1ZzUWdCb2FFdEVOMkJob216VHlUMlRXVEh5NW94RmQKY3d1T3FZM0hieUMwTTA3RXFEZWJEekVmeWpieEU2aXhGdUZxclVyd0xnZGN2T2JxcjFROS8zQUlyY3pWM1RKOQpNeCt0dUtTTldTWGZLaHA0eEFvVHRwVTFkQVJxTXRsNWtPNWVRUnNkSUpIYXdaY0JOWVRSbWpGNXM0T1M5cWZFCnl5NzNJUUtCZ0NuOS9GazBndmJhOStmZlFHV1daa2JqQnJvazJYVjFsZTQxWjUrcE5EclJBa2JycGtjaXRmN2MKNXdVenRnbXgyZkN4dWh6cCswWTdIZUM1eGNaVDBpUmJ3bGFIQUNhWGpJOVNYTjRzTFNCbGwzRG52dnRaR05FNQpoeW01c1MweC8vaUN6M1lEeC9hQ3JZaU00L09mc2Y1ODltdFBkbUd0dzZ3NUZWbWE2eklsCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==",  # noqa: E501 pylint: disable=line-too-long
    "AWS_THING_CERT": "LS0tLS2CRUdJTiBDRVJUSUZJQ0FURS0LLS0tCk1JSURXVENDQWtHZ0F3SUJBZ0lVTGlwelZiYWUweURjak5pWEpYdDZWV0t2TWJZd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1RURkxNRWtHQTFVRUN3eENRVzFoZW05dUlGZGxZaUJUWlhKMmFXTmxjeUJQUFVGdFlYcHZiaTVqYjIwZwpTVzVqTGlCTVBWTmxZWFIwYkdVZ1UxUTlWMkZ6YUdsdVozUnZiaUJEUFZWVE1CNFhEVEl3TURNeE1USXlNRGN3Ck1Gb1hEVFE1TVRJek1USXpOVGsxT1Zvd0hqRWNNQm9HQTFVRUF3d1RRVmRUSUVsdlZDQkRaWEowYVdacFkyRjAKWlRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTitObWlOZmlORllIZ0c5WVdSZgpmUGJxQXB1c0tuZUgwb3phNUFLeTJMamlyeHlrdlZwTy9YR0RzU21uemlzNThMa0dnZjlhYWhNc1JoSmhxV1RXCm1VRmVaRlVkOXNUbUpZamFlUS8vNHJpZm5Ea3kxQ1ZPMHVFbUtWTTJUa0FVSjJhOGNDUjJaQmFvMUpqOVNGMnIKUXh4aThUa2tVb1ZVUDlCam5sbzJDT01DSnlPTnVPOTRUY3lsSEdyWWJBU1h5TDdVdzZYWEhSRDNYQjF1VjhOQQpEVTVSSmgxRjVEN0duMGFjRElNYitCQ1Z2WGd0V29raFJIS3BWQlBJS2JNQUs3ejhiRGhIK0hBK1JZYnZwQXlyCkVUbjFBM25VdExLclN4V2FXa3VUTzVyY3NGazNweEgwNU1QWVBtcVM1T2hybElxMGZUZW5uTTRmNytORXJQVzcKMnhVQ0F3RUFBYU5nTUY0d0h3WURWUjBqQkJnd0ZvQVVxc0M5d0tHL0g5SVA5VUttY3pOTG9OZWtLcll3SFFZRApWUjBPQkJZRUZGWGM1cDR2QTZNU2thVUZ4T2tJc3ZWNS9CZ3pNQXdHQTFVZEV3RUIvd1FDTUFBd0RnWURWUjBQCkFRSC9CQVFEQWdlQU1BMEdDU3FHU0liM0RRRUJDd1VBQTRJQkFRQ1FZcGRwVWVxTXFRWktmNlo1OGQ1eFJuUFMKa3phR29oQk9NQ1BqREg4SmNQazFpei9IWmFPcnBtUE5ST0o0bzZWY0gvWWsyZ2RKVndwWHFhd28yRmJFbGhrWgpwYTlubnFtTnpGYzhmdHhPWVdsZFF3RWJiRmZEQTMwVUpUU3FWNW9mRXNoSTE1SnN6YUo3YTNkVGxMQTJPMVhHCnhPeEs1ZnZ6am8rMzVGYVJXdWxnTmtkMzdvM0dsUjUvdlF1NHY2QlVveGxYdUVoS2dMUnFHSDFxa2h5S2FtbEkKZzdPclVveTV5aGg5OHg5VStUU25OL0FaY3k2RTkwTVpHV0FEcVNUeVBnUXJCUmJYa2lNQU8vWEZIRitxWjhRVgpnSTVsdEFDbFhqdTY3M2phWmI5ZGVqZ29JK2IxTWRxVmRldXFjeHhsNWx5SFBnUUpobURNdVlLS1dBVWsKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=",  # noqa: E501 pylint: disable=line-too-long
    "AWS_ROOT_CERT": "LS0tLS2CRUdJTiBDRVJUSUZJQ0FURS0tqS0tCk1JSURRVENDQWltZ0F3SUJBZ0lUQm15Zno1bS9qQW81NHZCNGlrUG1salpieWpBTkJna3Foa2lHOXcwQkFRc0YKQURBNU1Rc3dDUVlEVlFRR0V3SlZVekVQTUEwR0ExVUVDaE1HUVcxaGVtOXVNUmt3RndZRFZRUURFeEJCYldGNgpiMjRnVW05dmRDQkRRU0F4TUI0WERURTFNRFV5TmpBd01EQXdNRm9YRFRNNE1ERXhOekF3TURBd01Gb3dPVEVMCk1Ba0dBMVVFQmhNQ1ZWTXhEekFOQmdOVkJBb1RCa0Z0WVhwdmJqRVpNQmNHQTFVRUF4TVFRVzFoZW05dUlGSnYKYjNRZ1EwRWdNVENDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFMSjRnSEhLZU5YagpjYTlIZ0ZCMGZXN1kxNGgyOUpsbzkxZ2hZUGwwaEFFdnJBSXRodE9nUTNwT3NxVFFOcm9Cdm8zYlNNZ0hGelpNCjlPNklJOGMrNnpmMXRSbjRTV2l3M3RlNWRqZ2RZWjZrL29JMnBlVktWdVJGNGZuOXRCYjZkTnFjbXpVNUwvcXcKSUZBR2JIclFnTEttK2Evc1J4bVBVRGdIM0tLSE9WajR1dFdwK1Vobk1KYnVsSGhlYjRtalVjQXdobWFoUldhNgpWT3VqdzVINVNOei8wZWd3TFgwdGRIQTExNGdrOTU3RVdXNjdjNGNYOGpKR0tMaEQrcmNkcXNxMDhwOGtEaTFMCjkzRmNYbW4vNnBVQ3l6aUtybEE0Yjl2N0xXSWJ4Y2NlVk9GMzRHZklENXlISTlZL1FDQi9JSURFZ0V3K095UW0KamdTdWJKcklxZzBDQXdFQUFhTkNNRUF3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFPQmdOVkhROEJBZjhFQkFNQwpBWVl3SFFZRFZSME9CQllFRklRWXpJVTA3THdNbEpRdUNGbWN4N0lRVGdvSU1BMEdDU3FHU0liM0RRRUJDd1VBCkE0SUJBUUNZOGpkYVFaQ2hHc1YyVVNnZ05pTU9ydVlvdTZyNGxLNUlwREIvRy93a2pVdTB5S0dYOXJieGVuREkKVTVQTUNDamptQ1hQSTZUNTNpSFRmSVVKclU2YWRUckNDMnFKZUhaRVJ4aGxiSTFCamp0L21zdjB0YWRRMXdVcwpOK2dEUzYzcFlhQUNidlh5OE1XeTdWdTMzUHFVWEhlZUU2Vi9VcTJWOHZpVE85NkxYRnZLV2xKYllLOFU5MHZ2Cm8vdWZRSlZ0TVZUOFF0UEhSaDhqcmRrUFNIQ2EyWFY0Y2RGeVF6UjFibGRad2dKY0ptQXB6eU1aRm82SVE2WFUKNU1zSSt5TVJRK2hES1hKaW9hbGRYZ2pVa0s2NDJNNFV3dEJWOG9iMnhKTkRkMlpod0xub1FkZVhlR0FEYmtweQpycVhSZmJvUW5vWnNHNHE1V1RQNDY4U1F2dkc1Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",  # noqa: E501 pylint: disable=line-too-long
}


@pytest.mark.parametrize("to_drop", _VARIABLES.keys())
def test_missing(monkeypatch: Any, tmp_path: Path, to_drop: str) -> None:
    """Run a test, dropping each of the keys in turn"""
    # Delete existing environment variables
    for var in AppConfig.variables():
        monkeypatch.delenv(var, raising=False)
    # Now set all variables, skipping variable under test
    for name, value in _VARIABLES.items():
        if name != to_drop:
            # Set the variable in the environment
            monkeypatch.setenv(name, str(value))
    # Always set the cert dir
    monkeypatch.setenv("CERT_DIR", str(tmp_path))
    # Expect an error
    with pytest.raises(ConfigError):
        AppConfig.from_env(dotenv=False)
