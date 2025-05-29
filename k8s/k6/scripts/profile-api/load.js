import http from "k6/http";
import { check, sleep } from "k6";
import { Rate, Trend, Counter } from "k6/metrics";
import { textSummary } from "k6/metrics";

// Custom metrics
const authErrorRate = new Rate("auth_errors");
const requestCounter = new Counter("total_requests");
const validationErrorRate = new Rate("validation_errors");
const connectionErrorRate = new Rate("connection_errors");
const duplicateEmailRate = new Rate("duplicate_email_errors");
const requestBodyErrorRate = new Rate("request_body_errors");

// Duration trends for each operation
const authTrend = new Trend("auth_duration");
const createTrend = new Trend("create_duration");
const updateTrend = new Trend("update_duration");
const deleteTrend = new Trend("delete_duration");
const listTrend = new Trend("list_duration");

// Response size metrics
const responseSizeTrend = new Trend("response_size");

// Test configuration
export const options = {
  stages: [
    { duration: "1m", target: 50 }, // Ramp up to 50 users
    { duration: "3m", target: 50 }, // Stay at 50 users
    { duration: "1m", target: 100 }, // Ramp up to 100 users
    { duration: "3m", target: 100 }, // Stay at 100 users
    { duration: "1m", target: 0 }, // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
    validation_errors: ["rate<0.01"], // Less than 1% validation errors
    connection_errors: ["rate<0.01"], // Less than 1% connection errors
    duplicate_email_errors: ["rate<0.01"], // Less than 1% duplicate email errors
    request_body_errors: ["rate<0.01"], // Less than 1% request body errors
  },
};

// Test data
const BASE_URL =
  __ENV.BASE_URL || "http://profile-api.default.svc.cluster.local";
const TEST_USER = {
  user_id: "FB",
  password: "FB.com",
};

// Helper function to get authentication token
function getAuthToken() {
  const loginPayload = JSON.stringify({
    user_id: TEST_USER.user_id,
    password: TEST_USER.password,
  });

  const loginResponse = http.post(
    `${BASE_URL}/api/v1/auth/token`,
    loginPayload,
    {
      headers: { "Content-Type": "application/json" },
    }
  );

  authTrend.add(loginResponse.timings.duration);
  responseSizeTrend.add(loginResponse.body.length);
  requestCounter.add(1);

  const checks = check(loginResponse, {
    "login status is 200": (r) => r.status === 200,
    "login returns token": (r) => {
      const data = r.json();
      return data && data.token !== undefined;
    },
  });

  if (!checks) {
    authErrorRate.add(1);
    console.error("Auth failed:", {
      status: loginResponse.status,
      body: loginResponse.body,
    });
    return null;
  }

  return loginResponse.json().token;
}

// Helper function to generate test data
function generateTestData() {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 1000);

  // Generate a random number of digits between 10 and 15
  const numDigits = Math.floor(Math.random() * 6) + 10; // Random number between 10 and 15
  // Generate a random number with exactly numDigits digits
  const phoneNumber = Math.floor(Math.random() * Math.pow(10, numDigits))
    .toString()
    .padStart(numDigits, "0");

  return {
    first_name: `Test${random}`,
    last_name: `User${random}`,
    email: `test${timestamp}${random}@example.com`,
    phone: `+1${phoneNumber}`,
    bio: `Test bio for user ${random}`,
    image_urls: [`https://example.com/images/test${random}.jpg`],
    address: {
      street: "123 Test St",
      city: "Test City",
      state: "TS",
      country: "Test Country",
      postal_code: "12345",
    },
  };
}

// Main test function
export default function () {
  // Get authentication token
  const token = getAuthToken();
  if (!token) {
    console.error("Failed to get authentication token");
    return;
  }

  const headers = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  // Test List Profiles
  const listResponse = http.get(`${BASE_URL}/api/v1/profiles`, { headers });
  listTrend.add(listResponse.timings.duration);

  const listChecks = check(listResponse, {
    "list profiles status is 200": (r) => r.status === 200,
    "list profiles has data": (r) => {
      const data = r.json();
      return Array.isArray(data);
    },
  });

  if (!listChecks) {
    console.error("List profiles failed:", {
      status: listResponse.status,
      body: listResponse.body,
    });
  }

  sleep(1);

  // Test Create Profile with valid data
  const testData = generateTestData();
  const createPayload = JSON.stringify(testData);

  const createResponse = http.post(
    `${BASE_URL}/api/v1/profiles`,
    createPayload,
    { headers }
  );
  createTrend.add(createResponse.timings.duration);

  const createChecks = check(createResponse, {
    "create profile status is 201": (r) => r.status === 201,
    "create profile has id": (r) => {
      const data = r.json();
      return data && data.profile && data.profile.id !== undefined;
    },
  });

  if (!createChecks) {
    if (createResponse.status === 400) {
      validationErrorRate.add(1);
    } else if (createResponse.status === 409) {
      duplicateEmailRate.add(1);
    } else if (createResponse.status === 413) {
      requestBodyErrorRate.add(1);
    } else if (createResponse.status >= 500) {
      connectionErrorRate.add(1);
    }
    console.error("Create profile failed:", {
      status: createResponse.status,
      body: createResponse.body,
    });
    return;
  }

  // Save the created profile ID
  const createdProfile = createResponse.json();
  const profileId = createdProfile.profile.id;

  sleep(1);

  // Test Get Profile
  const getResponse = http.get(`${BASE_URL}/api/v1/profiles/${profileId}`, {
    headers,
  });

  const getChecks = check(getResponse, {
    "get profile status is 200": (r) => r.status === 200,
    "get profile has correct data": (r) => {
      const data = r.json();
      return data && data.profile && data.profile.id === profileId;
    },
  });

  if (!getChecks) {
    console.error("Get profile failed:", {
      status: getResponse.status,
      body: getResponse.body,
    });
  }

  sleep(1);

  // Test Update Profile
  const updateData = generateTestData();
  const updatePayload = JSON.stringify(updateData);

  const updateResponse = http.put(
    `${BASE_URL}/api/v1/profiles/${profileId}`,
    updatePayload,
    { headers }
  );
  updateTrend.add(updateResponse.timings.duration);

  const updateChecks = check(updateResponse, {
    "update profile status is 200": (r) => r.status === 200,
    "update profile has updated data": (r) => {
      const data = r.json();
      return (
        data &&
        data.profile &&
        data.profile.first_name === updateData.first_name
      );
    },
  });

  if (!updateChecks) {
    if (updateResponse.status === 400) {
      validationErrorRate.add(1);
    } else if (updateResponse.status === 409) {
      duplicateEmailRate.add(1);
    } else if (updateResponse.status === 413) {
      requestBodyErrorRate.add(1);
    } else if (updateResponse.status >= 500) {
      connectionErrorRate.add(1);
    }
    console.error("Update profile failed:", {
      status: updateResponse.status,
      body: updateResponse.body,
    });
  }

  sleep(1);

  // Test Delete Profile
  const deleteResponse = http.del(
    `${BASE_URL}/api/v1/profiles/${profileId}`,
    null,
    { headers }
  );
  deleteTrend.add(deleteResponse.timings.duration);

  const deleteChecks = check(deleteResponse, {
    "delete profile status is 204": (r) => r.status === 204,
  });

  if (!deleteChecks) {
    console.error("Delete profile failed:", {
      status: deleteResponse.status,
      body: deleteResponse.body,
    });
  }
}

// Summary handler
export function handleSummary(data) {
  return {
    stdout: textSummary(data, { indent: " ", enableColors: true }),
    "summary.json": JSON.stringify(data),
  };
}
