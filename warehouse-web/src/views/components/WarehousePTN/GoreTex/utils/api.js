import axios from "axios";

const API_URL = (import.meta.env.VITE_API_URL || "/api/v1").replace(/\/$/, "");
const formsUrl = `${API_URL}/gore-tex/forms`;
const dashboardUrl = `${API_URL}/gore-tex/dashboard`;

export async function getGoreTexWeeklyDashboard(year, week) {
  const response = await axios.get(`${dashboardUrl}/weekly`, {
    params: { year, week },
  });
  return response.data.data;
}

export async function submitWaterproofForm(data, isEdit = false) {
  const response = await axios.post(`${formsUrl}/waterproof`, {
    line: data.meta.line,
    styleName: data.meta.styleName,
    inspectionDate: data.meta.inspectionDate,
    data,
    isEdit,
  });
  return response.data.data;
}

export async function submitCentrifugalForm(data, isEdit = false) {
  const response = await axios.post(`${formsUrl}/centrifugal`, {
    inspectionDate: data.inspectionDate,
    data,
    isEdit,
  });
  return response.data.data;
}

export async function submitAnalysisForm(data, analysisId = 0, isEdit = false) {
  const response = await axios.post(`${formsUrl}/analysis`, {
    analysisId,
    data,
    isEdit,
  });
  return response.data.data;
}

export async function getGoreTexForms() {
  const response = await axios.get(formsUrl);
  return response.data.data || [];
}

export async function getWaterproofForm(line, styleName) {
  const response = await axios.get(`${formsUrl}/waterproof`, {
    params: { line, styleName },
  });
  return response.data.data;
}

export async function getCentrifugalForm(inspectionDate) {
  const response = await axios.get(
    `${formsUrl}/centrifugal/${encodeURIComponent(inspectionDate)}`,
  );
  return response.data.data;
}

export async function getAnalysisForm(id) {
  const response = await axios.get(
    `${formsUrl}/analysis/${encodeURIComponent(id)}`,
  );
  return response.data.data;
}

export function getApiErrorMessage(error) {
  return (
    error?.response?.data?.message ||
    error?.message ||
    "Không thể kết nối tới máy chủ."
  );
}
