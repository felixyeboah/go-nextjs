"use client";

import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { AlertCircle } from "lucide-react";

export default function VerificationErrorContent() {
  const searchParams = useSearchParams();
  const error = searchParams.get("error") || "Verification failed";

  return (
    <>
      <div className="flex flex-col items-center space-y-4 text-center">
        <AlertCircle className="h-16 w-16 text-red-500" />
        <h1 className="text-2xl font-semibold tracking-tight">
          Verification Failed
        </h1>
        <p className="text-sm text-muted-foreground text-red-500">{error}</p>
      </div>
      <div className="flex flex-col space-y-4">
        <Link href="/auth/resend-verification">
          <Button className="w-full">Resend Verification Email</Button>
        </Link>
        <Link href="/auth/login">
          <Button variant="outline" className="w-full">
            Back to Login
          </Button>
        </Link>
      </div>
    </>
  );
}
