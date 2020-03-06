package genReports

const defaultGridTemplate = `
%use square brackets as golang text templating delimiters
\documentclass{article}
\usepackage{graphicx}
\usepackage[margin=0.5in]{geometry}

\graphicspath{ {images/} }
\begin{document}
\title{[[.Title]] [[if .VariableValues]] \\ \large [[.VariableValues]] [[end]] [[if .Description]] \\ \small [[.Description]] [[end]]}
\date{[[.FromFormatted]]\\to\\[[.ToFormatted]]}
\maketitle
\begin{center}
[[range .Panels]][[if .IsPartialWidth]]\begin{minipage}{[[.Width]]\textwidth}
\includegraphics[width=\textwidth]{image[[.Id]]}
\end{minipage}
[[else]]\par
\vspace{0.5cm}
\includegraphics[width=\textwidth]{image[[.Id]]}
\par
\vspace{0.5cm}
[[end]][[end]]

\end{center}
\end{document}
`
