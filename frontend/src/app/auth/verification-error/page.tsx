import { Suspense } from "react";
import VerificationErrorContent from "@/components/auth/verification-error-content";

export default function VerificationErrorPage() {
  return (
    <div className="container flex h-screen w-screen flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <Suspense fallback={<div>Loading...</div>}>
          <VerificationErrorContent />
        </Suspense>
      </div>
    </div>
  );
}
