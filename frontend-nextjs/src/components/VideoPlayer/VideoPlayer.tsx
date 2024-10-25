"use client";

import { useEffect, useRef } from "react";
import * as shaka from "shaka-player/dist/shaka-player.compiled";

export type VideoPlayerProps = {
  url: string;
  poster: string;
};

export default function VideoPlayer({ url, poster }: VideoPlayerProps) {
  const videoNodeRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<shaka.Player | null>(null);

  useEffect(() => {
    if (typeof window !== "undefined" && videoNodeRef.current && !playerRef.current) {
      shaka.polyfill.installAll();
      playerRef.current = new shaka.Player(videoNodeRef.current);
      playerRef.current
        .load(url)
        .catch(console.error);
    }

    return () => {
      playerRef.current?.destroy();
      playerRef.current = null;
    };
  }, [url]);

  return (
    <video
      ref={videoNodeRef}
      className="w-full h-full"
      controls
      autoPlay
      poster={poster}
    />
  );
}
