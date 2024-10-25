import { VideoModel } from "../../../models";

export async function getVideo(slug: string): Promise<VideoModel> {
  const response = await fetch(`http://localhost:8000/api/videos/${slug}`, {
    cache: "no-cache",
  });
  return response.json();
}

//revalidate on demand
// /admin/videos -> django -> http -> next.js -> revalidate

