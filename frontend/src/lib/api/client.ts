import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from "axios";

export class APIError extends Error {
  constructor(message: string, public status: number, public code?: string) {
    super(message);
    this.name = "APIError";
  }
}

export class APIClient {
  private client: AxiosInstance;

  constructor(
    baseURL: string = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
  ) {
    this.client = axios.create({
      baseURL,
      headers: {
        "Content-Type": "application/json",
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        // Get the token from localStorage or wherever you store it
        const token =
          typeof window !== "undefined" ? localStorage.getItem("token") : null;

        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response) {
          const status = error.response.status;
          const message =
            (error.response.data as any)?.message || error.message;
          const code = (error.response.data as any)?.code;

          // Handle 401 Unauthorized
          if (status === 401) {
            // Clear token and redirect to login
            if (typeof window !== "undefined") {
              localStorage.removeItem("token");
              window.location.href = "/auth/login";
            }
          }

          throw new APIError(message, status, code);
        }
        throw error;
      }
    );
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.get<T>(url, config);
    return response.data;
  }

  async post<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    const response = await this.client.post<T>(url, data, config);
    return response.data;
  }

  async put<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    const response = await this.client.put<T>(url, data, config);
    return response.data;
  }

  async patch<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    const response = await this.client.patch<T>(url, data, config);
    return response.data;
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.delete<T>(url, config);
    return response.data;
  }
}

// Create and export a singleton instance
export const api = new APIClient();
