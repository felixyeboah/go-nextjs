"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/use-toast";

export default function VerifyContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get("token");
  const { toast } = useToast();
  const [isLoading, setIsLoading] = useState(true);
  const [isVerified, setIsVerified] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!token) {
      setIsLoading(false);
      setError("Verification token is missing");
      return;
    }

    const verifyEmail = async () => {
      try {
        const response = await fetch(`/api/auth/verify?token=${token}`, {
          method: "GET",
        });

        if (!response.ok) {
          const data = await response.json();
          throw new Error(data.message || "Failed to verify email");
        }

        setIsVerified(true);
        toast({
          title: "Success",
          description: "Your email has been verified successfully",
        });
      } catch (error) {
        setError(
          error instanceof Error ? error.message : "Failed to verify email"
        );
        toast({
          title: "Error",
          description:
            error instanceof Error ? error.message : "Failed to verify email",
          variant: "destructive",
        });
      } finally {
        setIsLoading(false);
      }
    };

    verifyEmail();
  }, [token, toast]);

  return (
    <>
      <div className="flex flex-col space-y-2 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">
          Email Verification
        </h1>
        {isLoading ? (
          <p className="text-sm text-muted-foreground">
            Verifying your email...
          </p>
        ) : isVerified ? (
          <p className="text-sm text-muted-foreground">
            Your email has been verified successfully.
          </p>
        ) : (
          <p className="text-sm text-muted-foreground text-red-500">
            {error || "Failed to verify email"}
          </p>
        )}
      </div>

      {!isLoading && (
        <div className="flex flex-col space-y-4">
          {isVerified ? (
            <Button onClick={() => router.push("/auth/login")}>
              Go to Login
            </Button>
          ) : (
            <>
              <Button
                variant="outline"
                onClick={() => router.push("/auth/login")}
              >
                Go to Login
              </Button>
              <p className="px-8 text-center text-sm text-muted-foreground">
                <Link
                  href="/auth/resend-verification"
                  className="hover:text-brand underline underline-offset-4"
                >
                  Resend verification email
                </Link>
              </p>
            </>
          )}
        </div>
      )}
    </>
  );
}
