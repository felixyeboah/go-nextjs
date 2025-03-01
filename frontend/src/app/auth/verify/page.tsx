import { Suspense } from "react";
import VerifyContent from "@/components/auth/verify-content";

export default function VerifyPage() {
  return (
    <div className="container flex h-screen w-screen flex-col items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
        <Suspense fallback={<div>Loading...</div>}>
          <VerifyContent />
        </Suspense>
      </div>
    </div>
  );
}
