"use client";

import { Star } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { useEpisodeGradesQuery } from "../../hooks/use-episode-grades-query";
import { useMyEpisodeGradeQuery } from "../../hooks/use-my-episode-grade-query";
import { useSaveEpisodeGradeMutation } from "../../hooks/use-save-episode-grade-mutation";
import { GradeBar } from "../grade-bar";
import { GradeCard } from "../grade-card";
import { GradeSlider } from "../grade-slider";

interface EpisodeGradeSectionProps {
  episodeId: string;
  averageGrade: number | null | undefined;
  gradeCount: number | undefined;
}

export function EpisodeGradeSection({
  episodeId,
  averageGrade,
  gradeCount,
}: EpisodeGradeSectionProps) {
  const { t } = useTranslation();
  const { data: gradesData } = useEpisodeGradesQuery(episodeId);
  const { data: myGrade, isLoading: myGradeLoading } =
    useMyEpisodeGradeQuery(episodeId);

  if (myGradeLoading) {
    return (
      <Card>
        <CardHeader className="pb-2">
          <CardTitle className="flex items-center gap-2 text-sm font-medium text-gray-500 dark:text-gray-400">
            <Star className="h-4 w-4" />
            {t("episodeDetail.grades")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-24 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
        </CardContent>
      </Card>
    );
  }

  return (
    // key remounts Inner when myGrade transitions, so its useState re-seeds.
    <EpisodeGradeSectionInner
      key={myGrade ? "graded" : "ungraded"}
      episodeId={episodeId}
      averageGrade={averageGrade}
      gradeCount={gradeCount}
      items={gradesData?.grades ?? []}
      initialGrade={myGrade?.grade ?? 0.5}
      initialComment={myGrade?.comment ?? ""}
      hasExistingGrade={myGrade != null}
    />
  );
}

function EpisodeGradeSectionInner({
  episodeId,
  averageGrade,
  gradeCount,
  items,
  initialGrade,
  initialComment,
  hasExistingGrade,
}: {
  episodeId: string;
  averageGrade: number | null | undefined;
  gradeCount: number | undefined;
  items: ReadonlyArray<
    NonNullable<
      ReturnType<typeof useEpisodeGradesQuery>["data"]
    >["grades"][number]
  >;
  initialGrade: number;
  initialComment: string;
  hasExistingGrade: boolean;
}) {
  const { t } = useTranslation();
  const saveMutation = useSaveEpisodeGradeMutation();

  const [value, setValue] = useState(initialGrade);
  const [comment, setComment] = useState(initialComment);
  // Hide the slider until the user explicitly starts grading. Otherwise the
  // default 0.5 reads like an existing grade for someone who has not rated yet.
  const [isEditing, setIsEditing] = useState(hasExistingGrade);

  const handleSave = () => {
    saveMutation.mutate({
      episodeId,
      data: {
        grade: value,
        comment: comment.trim() === "" ? null : comment,
      },
    });
  };

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-sm font-medium text-gray-500 dark:text-gray-400">
          <Star className="h-4 w-4" />
          {t("episodeDetail.grades")}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <div>
          <p className="text-xs text-gray-500 dark:text-gray-400 mb-2">
            {t("episodeDetail.averageGrade")}
          </p>
          <GradeBar value={averageGrade ?? null} count={gradeCount} size="lg" />
        </div>

        {items.length === 0 ? (
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {t("episodeDetail.noGrades")}
          </p>
        ) : (
          <div className="space-y-2">
            {items.map((g) => (
              <GradeCard key={`${g.episode_id}-${g.user_id}`} grade={g} />
            ))}
          </div>
        )}

        <div className="space-y-3 border-t border-gray-200 dark:border-gray-700 pt-4">
          <p className="text-xs font-medium text-gray-700 dark:text-gray-300">
            {t("episodeDetail.yourGrade")}
          </p>
          {isEditing ? (
            <>
              <GradeSlider
                value={value}
                onChange={setValue}
                comment={comment}
                onCommentChange={setComment}
                disabled={saveMutation.isPending}
              />
              <div className="flex justify-end">
                <Button onClick={handleSave} disabled={saveMutation.isPending}>
                  {saveMutation.isPending
                    ? t("dialog.saving")
                    : t("dialog.save")}
                </Button>
              </div>
            </>
          ) : (
            <Button variant="outline" onClick={() => setIsEditing(true)}>
              <Star className="mr-2 h-4 w-4" />
              {t("episodeDetail.gradeThisEpisode")}
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
